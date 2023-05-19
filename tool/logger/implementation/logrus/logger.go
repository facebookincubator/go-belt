// Copyright 2022 Meta Platforms, Inc. and affiliates.
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

// Copyright (c) Facebook, Inc. and its affiliates.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package logrus

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/facebookincubator/go-belt"
	"github.com/facebookincubator/go-belt/pkg/field"
	"github.com/facebookincubator/go-belt/pkg/valuesparser"
	"github.com/facebookincubator/go-belt/tool/logger"
	"github.com/facebookincubator/go-belt/tool/logger/adapter"
	"github.com/facebookincubator/go-belt/tool/logger/types"
	"github.com/sirupsen/logrus"
)

var (
	// FieldNameTraceIDs is the field name used to store trace IDs.
	FieldNameTraceIDs = "trace_id"
)

// Emitter implements types.Emitter given a logrus entry.
type Emitter struct {
	// LogrusEntry is the actual emitter. And Emitter structure is only a wrapper
	// to implement another interface based on logrus.Entry.
	//
	// We use logrus.Entry instead of logrus.Logger to be able to store precompiled fields
	LogrusEntry *logrus.Entry
}

var _ types.Emitter = (*Emitter)(nil)

// Flush implements types.Emitter
func (l *Emitter) Flush() {}

// CheckLevel returns true if an event of a given logging level will be logged
func (l *Emitter) CheckLevel(level types.Level) bool {
	return LevelFromLogrus(l.LogrusEntry.Level) >= level
}

// Emit implements types.Emitter
func (l *Emitter) Emit(entry *types.Entry) {
	if entry.Level == types.LevelNone || (!l.CheckLevel(entry.Level) && entry.Level > types.LevelPanic) {
		return
	}
	l.logEntry(entry)
}

func ptr[T any](in T) *T {
	return &in
}

func (l *Emitter) logEntry(entry *types.Entry) {
	logrusLevel := LevelToLogrus(entry.Level)
	logrusEntry := ptr(*l.LogrusEntry) // shallow copy
	logrusEntry.Level = logrusLevel
	logrusEntry.Caller = nil
	if entry.Caller.Defined() {
		caller := &runtime.Frame{
			PC:   uintptr(entry.Caller),
			Func: runtime.FuncForPC(uintptr(entry.Caller)),
		}
		caller.Function = caller.Func.Name()
		caller.File, caller.Line = entry.Caller.FileLine()
		caller.Entry = entry.Caller.Entry()

		// this is effectively a no-op line, see RestoreCallerHook's description:
		logrusEntry.Caller = caller
		// the actualizator-line:
		logrusEntry.Context = context.WithValue(
			context.Background(), LogrusCtxKeyCaller, caller,
		)
	}
	logrusEntry.Time = entry.Timestamp
	if !entry.Properties.Has(EntryPropertyIgnoreFields) {
		newFieldsCount := 0
		if entry.Fields != nil {
			newFieldsCount = entry.Fields.Len()
		}
		if newFieldsCount > 0 || entry.TraceIDs != nil {
			fields := make(logrus.Fields, len(logrusEntry.Data)+newFieldsCount+1)
			for k, v := range logrusEntry.Data {
				fields[k] = v
			}
			if entry.Fields != nil {
				entry.Fields.ForEachField(func(f *field.Field) bool {
					fields[f.Key] = f.Value
					return true
				})
			}
			logrusEntry.Data = fields
			if entry.TraceIDs != nil {
				logrusEntry.Data[FieldNameTraceIDs] = entry.TraceIDs
			}
		}
	}
	logrusEntry.Log(logrusLevel, entry.Message)
}

// NewEmitter returns a new types.Emitter implementation based on a logrus logger.
// This functions takes ownership of the logrus Logger instance.
func NewEmitter(logger *logrus.Logger, logLevel logger.Level) *Emitter {
	// Enforcing Trace level, because we will take care of logging levels ourselves.
	//
	// This is required, because logrus does not support contextual logging levels,
	// so we handle them manually to provide with this feature.
	logger.Level = logrus.TraceLevel

	// Logrus internally duplicates Entry and does not copy the Caller, so
	// to keep the Caller we implement a workaround hook, which manually restores
	// it.
	//
	// Note:
	// The problem is there with any ReportCaller value, because with
	// ReportCaller disabled the Caller is just set to nil, but will
	// ReportCaller enabled it overwrites Caller using logrus's internal logic.
	// So either it is required to hack-around each Formatter or to use a Hook,
	// or to call `write` method manually (which is extra unsafety).
	logger.Hooks.Add(newRestoreCallerHook())
	logger.ReportCaller = false

	logrusEntry := newLogrusEntry(logger, LevelToLogrus(logLevel))
	return &Emitter{
		LogrusEntry: logrusEntry,
	}
}

func newLogrusEntry(logrusLogger *logrus.Logger, level logrus.Level) *logrus.Entry {
	return &logrus.Entry{
		Logger: logrusLogger,
		Level:  level,
	}
}

type mostlyPersistentData struct {
	fmtBufPool    *sync.Pool
	entryPool     *sync.Pool
	preHooks      types.PreHooks
	hooks         types.Hooks
	traceIDs      belt.TraceIDs
	getCallerFunc types.GetCallerPC
	messagePrefix string
}

// CompactLogger implements types.CompactLogger given a logrus logger.
type CompactLogger struct {
	*mostlyPersistentData
	emitter            *Emitter
	contextFields      *field.FieldsChain
	prepareEmitterOnce sync.Once
}

var _ adapter.CompactLogger = (*CompactLogger)(nil)

// Flush implements types.CompactLogger
func (l *CompactLogger) Flush() {
	l.emitter.Flush()
	for _, hook := range l.hooks {
		hook.Flush()
	}
}

var entryPropertiesIgnoreFields = types.EntryProperties{EntryPropertyIgnoreFields}

func (l *CompactLogger) acquireEntry() *types.Entry {
	entry := l.entryPool.Get().(*types.Entry)
	entry.Timestamp = time.Now()
	entry.Fields = l.contextFields
	entry.Properties = entryPropertiesIgnoreFields
	return entry
}

func (l *CompactLogger) releaseEntry(entry *types.Entry) {
	entry.Caller = 0
	entry.Fields = nil
	entry.Message = ""
	entry.Properties = nil
	l.entryPool.Put(entry)
}

func (l *CompactLogger) acquireBuf() *strings.Builder {
	return l.fmtBufPool.Get().(*strings.Builder)
}

func (l *CompactLogger) releaseBuf(buf *strings.Builder) {
	if buf.Cap() > 1024 {
		return
	}
	buf.Reset()
	l.fmtBufPool.Put(buf)
}

// LogEntry logs the given entry.
func (l *CompactLogger) LogEntry(entry *types.Entry) {
	if !l.emitter.CheckLevel(entry.Level) {
		return
	}
	l.logEntry(entry)
}
func (l *CompactLogger) logEntry(entry *types.Entry) {
	if !entry.Caller.Defined() && l.getCallerFunc != nil {
		entry.Caller = l.getCallerFunc()
	}
	if !adapter.ProcessHooks(l.hooks, entry) {
		return
	}
	l.prepareEmitter()
	l.emitter.logEntry(entry)
}

// WithMessagePrefix implements types.CompactLogger
func (l *CompactLogger) WithMessagePrefix(prefix string) adapter.CompactLogger {
	clone := l.clone()
	clone.messagePrefix += prefix
	return clone
}

// LogFields implements types.CompactLogger
func (l *CompactLogger) LogFields(level types.Level, message string, fields field.AbstractFields) {
	preHooksResult := adapter.LogFieldsPreprocess(l.preHooks, l.traceIDs, level, l.emitter.CheckLevel(level), message, fields)
	if preHooksResult.Skip {
		return
	}

	if preHooksResult.ExtraFields != nil {
		fields = field.Slice[field.AbstractFields]{fields, preHooksResult.ExtraFields}
	}
	logger := l
	if fields.Len() > 0 {
		logger = l.WithFields(fields).(*CompactLogger)
	}
	entry := logger.acquireEntry()
	defer logger.releaseEntry(entry)

	entry.Level = level
	entry.Message = message
	logger.logEntry(entry)
}

// Logf implements types.CompactLogger
func (l *CompactLogger) Logf(level types.Level, format string, args ...any) {
	preHooksResult := adapter.LogfPreprocess(l.preHooks, l.traceIDs, level, l.emitter.CheckLevel(level), format, args...)
	if preHooksResult.Skip {
		return
	}

	logger := l
	if preHooksResult.ExtraFields != nil {
		logger = logger.WithFields(preHooksResult.ExtraFields).(*CompactLogger)
	}

	entry := logger.acquireEntry()
	defer logger.releaseEntry(entry)

	buf := logger.acquireBuf()
	defer logger.releaseBuf(buf)
	fmt.Fprintf(buf, format, args...)

	entry.Level = level
	entry.Message = buf.String()
	logger.logEntry(entry)
}

// Log implements types.CompactLogger
func (l *CompactLogger) Log(level types.Level, values ...any) {
	preHooksResult := adapter.LogPreprocess(l.preHooks, l.traceIDs, level, l.emitter.CheckLevel(level), values...)
	if preHooksResult.Skip {
		return
	}

	if len(values) == 1 {
		if entry, ok := values[0].(*logger.Entry); ok {
			entry.Fields = field.Add(entry.Fields, preHooksResult.ExtraFields)
			if entry.TraceIDs != nil {
				entry.Fields = field.Add(entry.Fields, preHooksResult.ExtraFields, &field.Field{Key: FieldNameTraceIDs, Value: entry.TraceIDs})
			} else {
				entry.Fields = field.Add(entry.Fields, preHooksResult.ExtraFields)
			}
			l.logEntry(entry)
			return
		}
	}

	var finalFields field.AbstractFields
	valuesParser := valuesparser.AnySlice(values)
	if preHooksResult.ExtraFields != nil {
		finalFields = field.Slice[field.AbstractFields]{&valuesParser, preHooksResult.ExtraFields}
	} else {
		finalFields = &valuesParser
	}

	logger := l
	if finalFields.Len() != 0 {
		logger = logger.WithFields(finalFields).(*CompactLogger)
	}

	entry := logger.acquireEntry()
	defer logger.releaseEntry(entry)

	buf := logger.acquireBuf()
	defer logger.releaseBuf(buf)
	valuesParser.WriteUnparsed(buf)

	entry.Level = level
	entry.Message = buf.String()
	logger.logEntry(entry)
}

// Level implements types.CompactLogger
func (l *CompactLogger) Level() types.Level {
	return LevelFromLogrus(l.emitter.LogrusEntry.Level)
}

// WithLevel implements types.CompactLogger
func (l *CompactLogger) WithLevel(newLevel types.Level) adapter.CompactLogger {
	clone := l.clone()
	clone.emitter.LogrusEntry = newLogrusEntry(l.emitter.LogrusEntry.Logger, LevelToLogrus(newLevel))
	return clone
}

func (l *CompactLogger) branch() *CompactLogger {
	return &CompactLogger{
		mostlyPersistentData: l.mostlyPersistentData,
		emitter:              &Emitter{LogrusEntry: l.emitter.getLogrusEntry()},
		contextFields:        l.contextFields,
	}
}

func (l *CompactLogger) clone() *CompactLogger {
	clone := l.branch()
	clone.mostlyPersistentData = &[]mostlyPersistentData{*l.mostlyPersistentData}[0]
	return clone
}

// WithPreHooks implements types.CompactLogger
func (l *CompactLogger) WithPreHooks(preHooks ...types.PreHook) adapter.CompactLogger {
	clone := l.clone()
	clone.preHooks = make(types.PreHooks, len(l.preHooks)+len(preHooks))
	copy(clone.preHooks, l.preHooks)
	copy(clone.preHooks[len(l.preHooks):], preHooks)
	return clone
}

// WithHooks implements types.CompactLogger
func (l *CompactLogger) WithHooks(hooks ...types.Hook) adapter.CompactLogger {
	clone := l.clone()
	clone.hooks = make(types.Hooks, len(l.hooks)+len(hooks))
	copy(clone.hooks, l.hooks)
	copy(clone.hooks[len(l.hooks):], hooks)
	return clone
}

// WithField implements types.CompactLogger
func (l *CompactLogger) WithField(
	key field.Key,
	value field.Value,
	props ...field.Property,
) adapter.CompactLogger {
	branch := l.branch()
	branch.contextFields = l.contextFields.WithField(key, value, props)
	return branch
}

// WithFields implements types.CompactLogger
func (l *CompactLogger) WithFields(fields field.AbstractFields) adapter.CompactLogger {
	branch := l.branch()
	branch.contextFields = l.contextFields.WithFields(fields)
	return branch
}

// WithTraceIDs implements types.CompactLogger
func (l *CompactLogger) WithTraceIDs(allTraceIDs belt.TraceIDs, newTraceIDsCount int) belt.Tool {
	clone := l.clone()
	clone.traceIDs = allTraceIDs
	clone.contextFields = clone.contextFields.WithField(FieldNameTraceIDs, allTraceIDs)
	return clone
}

// WithContextFields implements types.Tool
func (l *CompactLogger) WithContextFields(allFields *field.FieldsChain, newFieldsCount int) belt.Tool {
	branch := l.branch()
	branch.contextFields = allFields
	return branch
}

func (l *CompactLogger) compileEmitterFields() {
	m := make(logrus.Fields, l.contextFields.Len())
	l.contextFields.ForEachField(func(f *field.Field) bool {
		m[f.Key] = f.Value
		return true
	})
	old := l.emitter.LogrusEntry
	l.emitter.setLogrusEntry(&logrus.Entry{
		Logger:  old.Logger,
		Data:    m,
		Time:    old.Time,
		Level:   old.Level,
		Caller:  old.Caller,
		Message: old.Message,
		Context: old.Context,
	})
}

func (l *CompactLogger) prepareEmitter() {
	l.prepareEmitterOnce.Do(func() {
		l.compileEmitterFields()
	})
}

// Emitter implements types.Emitter
func (l *CompactLogger) Emitter() types.Emitter {
	l.prepareEmitter()
	return l.emitter
}

func newCompactLoggerFromLogrus(logrusLogger *logrus.Logger, level logger.Level, opts ...types.Option) *CompactLogger {
	cfg := types.Options(opts).Config()
	return &CompactLogger{
		emitter: NewEmitter(logrusLogger, level),
		mostlyPersistentData: &mostlyPersistentData{
			getCallerFunc: cfg.GetCallerFunc,
			fmtBufPool: &sync.Pool{
				New: func() any {
					return &strings.Builder{}
				},
			},
			entryPool: &sync.Pool{
				New: func() any {
					return &types.Entry{}
				},
			},
		},
	}
}

// New returns a types.Logger implementation based on a logrus Logger.
//
// This function takes ownership of the logger instance.
func New(logger *logrus.Logger, opts ...types.Option) types.Logger {
	return adapter.GenericSugar{
		CompactLogger: newCompactLoggerFromLogrus(logger, LevelFromLogrus(logger.Level), opts...),
	}
}

// DefaultLogrusLogger returns a logrus logger with the default configuration.
//
// The configuration might be changed in future.
//
// Overwritable.
//
// Do not override this anywhere but in the `main` package.
var DefaultLogrusLogger = func() *logrus.Logger {
	l := logrus.New()
	l.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(caller *runtime.Frame) (string, string) {
			return "", fmt.Sprintf("%s:%d", filepath.Base(caller.File), caller.Line)
		},
	}
	return l
}

// Default returns a types.Logger based on logrus, using the default configuration.
//
// The configuration might be changed in future. Use this function if you would like
// to delegate logger configuration to somebody else and you do not rely on specific
// output format.
func Default() types.Logger {
	return New(DefaultLogrusLogger())
}
