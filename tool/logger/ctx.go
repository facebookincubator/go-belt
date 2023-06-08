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

package logger

import (
	"context"

	"github.com/facebookincubator/go-belt"
	"github.com/facebookincubator/go-belt/pkg/field"
)

// FromCtx returns the Logger defined in the context. If one is not defined,
// then the default Logger is returned (see function "Default").
func FromCtx(ctx context.Context) Logger {
	return FromBelt(belt.CtxBelt(ctx))
}

// CtxWithLogger returns a context with the given Logger set.
func CtxWithLogger(ctx context.Context, logger Logger) context.Context {
	return belt.WithTool(ctx, ToolID, logger)
}

// Flush forces to flush all buffers.
func Flush(ctx context.Context) {
	FromCtx(ctx).Flush()
}

// WithContextFields is a shorthand for FromCtx(ctx).WithContextFields
func WithContextFields(ctx context.Context, allFields *field.FieldsChain, newFieldsCount int) Logger {
	tool := FromCtx(ctx).WithContextFields(allFields, newFieldsCount)
	if l, ok := tool.(Logger); ok {
		return l
	}
	return Default()
}

// WithTraceIDs is a shorthand for FromCtx(ctx).WithTraceIDs.
func WithTraceIDs(ctx context.Context, traceIDs belt.TraceIDs, newTraceIDsCount int) Logger {
	tool := FromCtx(ctx).WithTraceIDs(traceIDs, newTraceIDsCount)
	if l, ok := tool.(Logger); ok {
		return l
	}
	return Default()
}

// Logf logs an unstructured message. All contextual structured
// fields are also logged.
//
// This method exists mostly for convenience, for people who
// has not got used to proper structured logging, yet.
// See `LogFields` and `Log`. If one have variables they want to
// log, it is better for scalable observability to log them
// as structured values, instead of injecting them into a
// non-structured string.
func Logf(ctx context.Context, level Level, format string, args ...any) {
	FromCtx(ctx).Logf(level, format, args...)
}

// LogFields logs structured fields with a explanation message.
//
// Anything that implements field.AbstractFields might be used
// as a collection of fields to be logged.
//
// Examples:
//
//	logger.LogFields(ctx, logger.LevelDebug, "new_request", field.Fields{{Key: "user_id", Value: userID}, {Key: "group_id", Value: groupID}})
//	logger.LogFields(ctx, logger.LevelInfo, "affected entries", field.Field{Key: "mysql_affected", Value: affectedRows})
//	logger.LogFields(ctx, logger.LevelError, "unable to fetch user info", request) // where `request` implements field.AbstractFields
//
// Sometimes it is inconvenient to manually describe each field,
// and for such cases see method `Log`.
func LogFields(ctx context.Context, level Level, message string, fields field.AbstractFields) {
	FromCtx(ctx).LogFields(level, message, fields)
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
//		logger.Log(ctx, logger.LevelError, err)
//		return err
//	}
//	logger.Log(ctx, logger.LevelDebug, "current user", user) // fields of structure "user" will be logged
//	logger.Log(ctx, logger.LevelDebug, map[string]any{"user_id": user.ID, "group_id", user.GroupID})
//	logger.Log(ctx, logger.LevelDebug, field.Fields{{Key: "user_id", Value: user.ID}, {Key: "group_id", Value: user.GroupID}})
//	logger.Log(ctx, logger.LevelDebug, "current user ID is ", user.ID, " and group ID is ", user.GroupID) // will result into message "current user ID is 1234 and group ID is 5678".
func Log(ctx context.Context, level Level, values ...any) {
	FromCtx(ctx).Log(level, values...)
}

// GetEmitter returns the Emitter (see the description of interface "Emitter").
func GetEmitter(ctx context.Context) Emitter {
	return FromCtx(ctx).Emitter()
}

// GetLevel returns the current logging level (see description of "Level").
func GetLevel(ctx context.Context) Level {
	return FromCtx(ctx).Level()
}

// TraceFields is just a shorthand for LogFields(ctx, logger.LevelTrace, ...)
func TraceFields(ctx context.Context, message string, fields field.AbstractFields) {
	FromCtx(ctx).TraceFields(message, fields)
}

// DebugFields is just a shorthand for LogFields(ctx, logger.LevelDebug, ...)
func DebugFields(ctx context.Context, message string, fields field.AbstractFields) {
	FromCtx(ctx).DebugFields(message, fields)
}

// InfoFields is just a shorthand for LogFields(ctx, logger.LevelInfo, ...)
func InfoFields(ctx context.Context, message string, fields field.AbstractFields) {
	FromCtx(ctx).InfoFields(message, fields)
}

// WarnFields is just a shorthand for LogFields(ctx, logger.LevelWarn, ...)
func WarnFields(ctx context.Context, message string, fields field.AbstractFields) {
	FromCtx(ctx).WarnFields(message, fields)
}

// ErrorFields is just a shorthand for LogFields(ctx, logger.LevelError, ...)
func ErrorFields(ctx context.Context, message string, fields field.AbstractFields) {
	FromCtx(ctx).ErrorFields(message, fields)
}

// PanicFields is just a shorthand for LogFields(ctx, logger.LevelPanic, ...)
//
// Be aware: Panic level also triggers a `panic`.
func PanicFields(ctx context.Context, message string, fields field.AbstractFields) {
	FromCtx(ctx).PanicFields(message, fields)
}

// FatalFields is just a shorthand for LogFields(ctx, logger.LevelFatal, ...)
//
// Be aware: Panic level also triggers an `os.Exit`.
func FatalFields(ctx context.Context, message string, fields field.AbstractFields) {
	FromCtx(ctx).FatalFields(message, fields)
}

// Trace is just a shorthand for Log(ctx, logger.LevelTrace, ...)
func Trace(ctx context.Context, values ...any) {
	FromCtx(ctx).Trace(values...)
}

// Debug is just a shorthand for Log(ctx, logger.LevelDebug, ...)
func Debug(ctx context.Context, values ...any) {
	FromCtx(ctx).Debug(values...)
}

// Info is just a shorthand for Log(ctx, logger.LevelInfo, ...)
func Info(ctx context.Context, values ...any) {
	FromCtx(ctx).Info(values...)
}

// Warn is just a shorthand for Log(ctx, logger.LevelWarn, ...)
func Warn(ctx context.Context, values ...any) {
	FromCtx(ctx).Warn(values...)
}

// Error is just a shorthand for Log(ctx, logger.LevelError, ...)
func Error(ctx context.Context, values ...any) {
	FromCtx(ctx).Error(values...)
}

// Panic is just a shorthand for Log(ctx, logger.LevelPanic, ...)
//
// Be aware: Panic level also triggers a `panic`.
func Panic(ctx context.Context, values ...any) {
	FromCtx(ctx).Panic(values...)
}

// Fatal is just a shorthand for Log(logger.LevelFatal, ...)
//
// Be aware: Fatal level also triggers an `os.Exit`.
func Fatal(ctx context.Context, values ...any) {
	FromCtx(ctx).Fatal(values...)
}

// Tracef is just a shorthand for Logf(ctx, logger.LevelTrace, ...)
func Tracef(ctx context.Context, format string, args ...any) {
	FromCtx(ctx).Tracef(format, args...)
}

// Debugf is just a shorthand for Logf(ctx, logger.LevelDebug, ...)
func Debugf(ctx context.Context, format string, args ...any) {
	FromCtx(ctx).Debugf(format, args...)
}

// Infof is just a shorthand for Logf(ctx, logger.LevelInfo, ...)
func Infof(ctx context.Context, format string, args ...any) {
	FromCtx(ctx).Infof(format, args...)
}

// Warnf is just a shorthand for Logf(ctx, logger.LevelWarn, ...)
func Warnf(ctx context.Context, format string, args ...any) {
	FromCtx(ctx).Warnf(format, args...)
}

// Errorf is just a shorthand for Logf(ctx, logger.LevelError, ...)
func Errorf(ctx context.Context, format string, args ...any) {
	FromCtx(ctx).Errorf(format, args...)
}

// Panicf is just a shorthand for Logf(ctx, logger.LevelPanic, ...)
//
// Be aware: Panic level also triggers a `panic`.
func Panicf(ctx context.Context, format string, args ...any) {
	FromCtx(ctx).Panicf(format, args...)
}

// Fatalf is just a shorthand for Logf(ctx, logger.LevelFatal, ...)
//
// Be aware: Fatal level also triggers an `os.Exit`.
func Fatalf(ctx context.Context, format string, args ...any) {
	FromCtx(ctx).Fatalf(format, args...)
}

// WithLevel returns a logger with logger level set to the given argument.
//
// See also the description of type "Level".
func WithLevel(ctx context.Context, level Level) Logger {
	return FromCtx(ctx).WithLevel(level)
}

// WithPreHooks returns a Logger which includes/appends pre-hooks from the arguments.
//
// See also description of "PreHook".
//
// Special case: to reset hooks use `WithPreHooks()` (without any arguments).
func WithPreHooks(ctx context.Context, preHooks ...PreHook) Logger {
	return FromCtx(ctx).WithPreHooks(preHooks...)
}

// WithHooks returns a Logger which includes/appends hooks from the arguments.
//
// See also description of "Hook".
//
// Special case: to reset hooks use `WithHooks()` (without any arguments).
func WithHooks(ctx context.Context, hooks ...Hook) Logger {
	return FromCtx(ctx).WithHooks(hooks...)
}

// WithField returns the logger with the added field (used for structured logging).
func WithField(ctx context.Context, key string, value any, props ...field.Property) Logger {
	return FromCtx(ctx).WithField(key, value, props...)
}

// WithFields returns the logger with the added fields (used for structured logging)
func WithFields(ctx context.Context, fields field.AbstractFields) Logger {
	return FromCtx(ctx).WithFields(fields)
}

// WithMessagePrefix adds a string to all messages logged through the derived logger.
func WithMessagePrefix(ctx context.Context, prefix string) Logger {
	return FromCtx(ctx).WithMessagePrefix(prefix)
}
