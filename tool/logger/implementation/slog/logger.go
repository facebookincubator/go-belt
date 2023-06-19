// Copyright 2023 Meta Platforms, Inc. and affiliates.
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package slog

import (
	"bytes"
	"context"
	"fmt"

	"github.com/facebookincubator/go-belt"
	"github.com/facebookincubator/go-belt/pkg/field"
	"github.com/facebookincubator/go-belt/pkg/valuesparser"
	"github.com/facebookincubator/go-belt/tool/logger/adapter"
	"github.com/facebookincubator/go-belt/tool/logger/implementation/stdlib"
	"github.com/facebookincubator/go-belt/tool/logger/types"
	"github.com/go-ng/xsync"
	"golang.org/x/exp/slog"
)

var (
	// FieldNameTraceIDs is the field name used to store belt.TraceIDs.
	FieldNameTraceIDs = "trace_id"

	// InternalErrorLogger is the logger used to print internal errors of this logger.
	InternalErrorLogger = stdlib.Default()
)

type pooledData struct {
	Entry types.Entry
	Buf   bytes.Buffer
}

// ConvertedLogger is an implementation of logger.Logger given an slog.Logger.
type ConvertedLogger struct {
	*Logger

	level           types.Level
	fields          *field.FieldsChain
	contextFields   *field.FieldsChain
	getCallerFunc   types.GetCallerPC
	pool            *xsync.PoolR[pooledData]
	traceIDs        belt.TraceIDs
	preHooks        types.PreHooks
	hooks           types.Hooks
	messagePrefix   string
	entryProperties types.EntryProperties
}

func newConvertedLogger(l *slog.Logger) *ConvertedLogger {
	return &ConvertedLogger{
		Logger: l,
		pool: xsync.NewPoolR(
			nil,
			func(t *pooledData) {
				t.Buf.Reset()
				t.Entry.Caller = 0
				t.Entry.Fields = nil
				t.Entry.Message = ""
				t.Entry.TraceIDs = nil
				t.Entry.Properties = t.Entry.Properties[:0]
			},
		),
	}
}

// New returns an instance of Logger given an slog Logger.
func New(l *slog.Logger) types.Logger {
	return adapter.LoggerFromCompactLogger(newConvertedLogger(l))
}

// NewFromHandler returns an instance of Logger given an slog Handler.
func NewFromHandler(h slog.Handler) types.Logger {
	return adapter.LoggerFromCompactLogger(newConvertedLogger(slog.New(h)))
}

// Default is the overridable function to return the default Logger of this implementation.
var Default = func() types.Logger {
	return New(slog.Default())
}

var _ adapter.CompactLogger = (*ConvertedLogger)(nil)

func (l *ConvertedLogger) checkLevel(level types.Level) bool {
	return l.Logger.Handler().Enabled(context.TODO(), LevelToSlog(level))
}

func (l *ConvertedLogger) logEntry(entry *types.Entry) {
	if entry.TraceIDs == nil {
		entry.TraceIDs = l.traceIDs
	}
	if !entry.Caller.Defined() && l.getCallerFunc != nil {
		entry.Caller = l.getCallerFunc()
	}
	if len(l.entryProperties) > 0 {
		entry.Properties = append(entry.Properties, l.entryProperties)
	}
	fields := entry.Fields
	entry.Fields = field.Add(fields, l.fields, l.contextFields)
	if !adapter.ProcessHooks(l.hooks, entry) {
		return
	}

	// TODO: fix the context (see https://github.com/facebookincubator/go-belt/issues/32)
	ctx := context.TODO()

	if !l.Enabled(ctx, LevelToSlog(entry.Level)) {
		return
	}

	r := slog.NewRecord(
		entry.Timestamp,
		LevelToSlog(entry.Level),
		entry.Message,
		uintptr(entry.Caller),
	)
	if len(entry.TraceIDs) != 0 {
		r.AddAttrs(slog.Attr{
			Key:   FieldNameTraceIDs,
			Value: slog.AnyValue(entry.TraceIDs),
		})
	}
	fieldsToAttrs(&r, fields) // we pass `fields` instead of `entry.Fields`, because `Handler` already received the other fields in WithContextFields and WithField(s).
	err := l.Logger.Handler().Handle(ctx, r)
	if err != nil {
		InternalErrorLogger.Error(err)
	}
}

// Logf logs an unstructured message. Though, of course, all
// contextual structured fields will also be logged.
//
// This method exists mostly for convenience, for people who
// has not got used to proper structured logging, yet.
// See `LogFields` and `Log`. If one have variables they want to
// log, it is better for scalable observability to log them
// as structured values, instead of injecting them into a
// non-structured string.
func (l *ConvertedLogger) Logf(level types.Level, format string, args ...any) {
	preHooksResult := adapter.LogfPreprocess(l.preHooks, l.traceIDs, level, l.checkLevel(level), format, args...)
	if preHooksResult.Skip {
		return
	}

	_pooledData := l.pool.Get()
	defer _pooledData.Release()
	entry := &_pooledData.Value.Entry
	buf := &_pooledData.Value.Buf

	buf.WriteString(l.messagePrefix)
	fmt.Fprintf(buf, format, args...)
	entry.Message = buf.String()

	entry.Level = level
	entry.Fields = field.Add(preHooksResult.ExtraFields)

	l.logEntry(entry)
}

// LogFields logs structured fields with a explanation message.
//
// Anything that implements field.AbstractFields might be used
// as a collection of fields to be logged.
//
// Examples:
//
//	l.LogFields(logger.LevelDebug, "new_request", field.Fields{{Key: "user_id", Value: userID}, {Key: "group_id", Value: groupID}})
//	l.LogFields(logger.LevelInfo, "affected entries", field.Field{Key: "mysql_affected", Value: affectedRows})
//	l.LogFields(logger.LevelError, "unable to fetch user info", request) // where `request` implements field.AbstractFields
//
// Sometimes it is inconvenient to manually describe each field,
// and for such cases see method `Log`.
func (l *ConvertedLogger) LogFields(level types.Level, message string, fields field.AbstractFields) {
	preHooksResult := adapter.LogFieldsPreprocess(l.preHooks, l.traceIDs, level, l.checkLevel(level), message, fields)
	if preHooksResult.Skip {
		return
	}

	_pooledData := l.pool.Get()
	defer _pooledData.Release()
	entry := &_pooledData.Value.Entry
	buf := &_pooledData.Value.Buf

	entry.Level = level
	entry.Fields = field.Add(fields, preHooksResult.ExtraFields)

	buf.WriteString(l.messagePrefix)
	buf.WriteString(message)
	entry.Message = buf.String()
	l.logEntry(entry)
}

// Log extracts structured fields from provided values, joins
// the rest into an unstructured message and logs the result.
//
// This function provides convenience (relatively to LogFields)
// at cost of a bit of performance.
//
// There are few ways to extract structured fields, which are
// applied for each value from `values` (in descending priority order):
//  1. If a `value` is an `*Entry` then the Entry is used (with its fields).
//     This works only if this is the only argument. Otherwise it is
//     threated as a simple structure (see point #3).
//  2. If a `value` implements field.AbstractFields then ForEachField method
//     is used (so it is become similar to LogFields).
//  3. If a `value` is a structure (or a pointer to a structure) then
//     fields of the structure are interpreted as structured fields
//     to be logged (see explanation below).
//  4. If a `value` is a map then fields a constructed out of this map.
//
// Structured arguments are processed effectively the same
// as they would have through a sequence of WithField/WithFields.
//
// Everything that does not fit into any of the rules above is just
// joined into an unstructured message (and works the same way
// as `message` in LogFields).
//
// How structures are parsed:
// Structures are parsed recursively. Each field name of the path in a tree
// of structures is added to the resulting field name (for example in "struct{A struct{B int}}"
// the field name will be `A.B`).
// To enforce another name use tag `log` (for example "struct{A int `log:"anotherName"`}"),
// to prevent a field from logging use tag `log:"-"`.
//
// Examples:
//
//	user, err := getUser()
//	if err != nil {
//		l.Log(logger.LevelError, err)
//		return err
//	}
//	l.Log(logger.LevelDebug, "current user", user) // fields of structure "user" will be logged
//	l.Log(logger.LevelDebug, map[string]any{"user_id": user.ID, "group_id", user.GroupID})
//	l.Log(logger.LevelDebug, field.Fields{{Key: "user_id", Value: user.ID}, {Key: "group_id", Value: user.GroupID}})
//	l.Log(logger.LevelDebug, "current user ID is ", user.ID, " and group ID is ", user.GroupID) // will result into message "current user ID is 1234 and group ID is 5678".
func (l *ConvertedLogger) Log(level types.Level, values ...any) {
	forceProcess := level == types.LevelFatal || level == types.LevelPanic

	preHooksResult := adapter.LogPreprocess(l.preHooks, l.traceIDs, level, l.checkLevel(level), values...)
	if preHooksResult.Skip && !forceProcess {
		return
	}

	if len(values) == 1 {
		if entry, ok := values[0].(*types.Entry); ok {
			l.logEntry(entry)
			return
		}
	}

	valuesParser := valuesparser.AnySlice(values)

	_pooledData := l.pool.Get()
	defer _pooledData.Release()
	entry := &_pooledData.Value.Entry
	buf := &_pooledData.Value.Buf

	buf.WriteString(l.messagePrefix)
	valuesParser.ExtractUnparsed(buf)
	entry.Message = buf.String()

	entry.Level = level
	entry.Fields = field.Add(valuesParser, preHooksResult.ExtraFields)

	l.logEntry(entry)
}

// Emitter returns the Emitter (see the description of interface "Emitter").
func (l *ConvertedLogger) Emitter() types.Emitter {
	return Emitter{Handler: l.Logger.Handler()}
}

// Level returns the current logging level (see description of "Level").
func (l *ConvertedLogger) Level() types.Level {
	return l.level
}

func (l ConvertedLogger) clone() *ConvertedLogger {
	return &l
}

// WithLevel returns a logger with logger level set to the given argument.
//
// See also the description of type "Level".
func (l *ConvertedLogger) WithLevel(level types.Level) adapter.CompactLogger {
	l = l.clone()
	l.level = level
	return l
}

// WithPreHooks returns a Logger which includes/appends pre-hooks from the arguments.
//
// See also description of "PreHook".
//
// Special case: to reset hooks use `WithPreHooks()` (without any arguments).
func (l *ConvertedLogger) WithPreHooks(preHooks ...types.PreHook) adapter.CompactLogger {
	l.clone()
	newPreHooks := make(types.PreHooks, len(l.preHooks)+len(preHooks))
	copy(newPreHooks, l.preHooks)
	copy(newPreHooks[len(l.preHooks):], preHooks)
	l.preHooks = newPreHooks
	return l
}

// WithHooks returns a Logger which includes/appends hooks from the arguments.
//
// See also description of "Hook".
//
// Special case: to reset hooks use `WithHooks()` (without any arguments).
func (l *ConvertedLogger) WithHooks(hooks ...types.Hook) adapter.CompactLogger {
	l.clone()
	newHooks := make(types.Hooks, len(l.hooks)+len(hooks))
	copy(newHooks, l.hooks)
	copy(newHooks[len(l.hooks):], hooks)
	l.hooks = newHooks
	return l
}

// WithField returns the logger with the added field (used for structured logging).
func (l *ConvertedLogger) WithField(key string, value any, props ...field.Property) adapter.CompactLogger {
	l = l.clone()
	l.Logger = slog.New(l.Logger.Handler().WithAttrs([]slog.Attr{{
		Key:   key,
		Value: slog.AnyValue(value),
	}}))
	l.fields = l.fields.WithField(key, value, props...)
	return l
}

// WithFields returns the logger with the added fields (used for structured logging)
func (l *ConvertedLogger) WithFields(fields field.AbstractFields) adapter.CompactLogger {
	l = l.clone()
	h := l.Logger.Handler()
	fields.ForEachField(func(f *field.Field) bool {
		h = h.WithAttrs([]slog.Attr{{
			Key:   f.Key,
			Value: slog.AnyValue(f.Value),
		}})
		return true
	})
	l.Logger = slog.New(h)
	l.fields = l.fields.WithFields(fields)
	return l
}

// WithMessagePrefix adds a string to all messages logged through the derived logger.
func (l *ConvertedLogger) WithMessagePrefix(prefix string) adapter.CompactLogger {
	l.clone()
	l.messagePrefix = prefix + l.messagePrefix
	return l
}

// Flush forces to flush all buffers.
func (l *ConvertedLogger) Flush() {}

// WithContextFields sets new context-defined fields. Supposed to be called
// only by an Belt.
//
// allFields contains all fields as a chain of additions in a reverse-chronological order,
// while newFieldsCount tells about how much of the fields are new (since last
// call of WithContextFields). Thus if one will call
// field.Slice(allFields, 0, newFieldsCount) they will get only the new fields.
// At the same time some Tool-s may prefer just to re-set all the fields instead of adding
// only new fields (due to performance reasons) and they may just use `allFields`.
func (l *ConvertedLogger) WithContextFields(allFields *field.FieldsChain, newFieldsCount int) belt.Tool {
	l = l.clone()
	l.contextFields = allFields
	attrs := make([]slog.Attr, 0, newFieldsCount)
	count := 0
	allFields.ForEachField(func(f *field.Field) bool {
		if count >= newFieldsCount {
			return false
		}
		attrs = append(attrs, fieldToAttr(f))
		return true
	})
	l.Logger = slog.New(l.Logger.Handler().WithAttrs(attrs))
	return l
}

// WithTraceIDs sets new context-defined TraceIDs. Supposed to be called
// only by an Belt.
//
// traceIDs and newTraceIDsCount has similar properties as allFields and newFieldsCount
// in the WithContextFields method.
func (l *ConvertedLogger) WithTraceIDs(traceIDs belt.TraceIDs, newTraceIDsCount int) belt.Tool {
	l = l.clone()
	l.traceIDs = traceIDs
	return l
}

// WithEntryProperties adds props to EntryProperties of each emitted Entry.
// This could be used only for enabling implementation-specific behavior.
func (l *ConvertedLogger) WithEntryProperties(props ...types.EntryProperty) adapter.CompactLogger {
	l = l.clone()
	l.entryProperties = l.entryProperties.Add(props...)
	return l
}

// Emitter is an implementation of logger.Emitter given an slog.Handler.
type Emitter struct {
	Handler
}

var _ types.Emitter = (*Emitter)(nil)

// Flush forces to flush all buffers.
func (Emitter) Flush() {}

// Emit just logs the provided Entry. It does not modify it.
func (e Emitter) Emit(entry *types.Entry) {
	// TODO: fix the context (see https://github.com/facebookincubator/go-belt/issues/32)
	ctx := context.TODO()

	r := slog.NewRecord(
		entry.Timestamp,
		LevelToSlog(entry.Level),
		entry.Message,
		uintptr(entry.Caller),
	)
	if len(entry.TraceIDs) != 0 {
		r.AddAttrs(slog.Attr{
			Key:   FieldNameTraceIDs,
			Value: slog.AnyValue(entry.TraceIDs),
		})
	}
	fieldsToAttrs(&r, entry.Fields)
	err := e.Handler.Handle(ctx, r)
	if err != nil {
		InternalErrorLogger.Error(err)
	}
}

func (e Emitter) Level() types.Level {
	for level := types.LevelTrace; level >= types.LevelFatal; level-- {
		if e.Handler.Enabled(context.TODO(), LevelToSlog(level)) {
			return level
		}
	}

	return types.LevelUndefined
}
