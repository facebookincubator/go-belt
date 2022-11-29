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

package adapter

import (
	"github.com/facebookincubator/go-belt"
	"github.com/facebookincubator/go-belt/pkg/field"
	"github.com/facebookincubator/go-belt/tool/logger/types"
)

// GenericSugar is a generic implementation of sugar-methods for a given types.CompactLogger.
//
// It does not contain any logic, it is just a bunch of method-wrappers to provide
// logger.Logger implementation using a given types.CompactLogger.
type GenericSugar struct {
	CompactLogger CompactLogger
}

var _ types.Logger = (*GenericSugar)(nil)

// Level implements logger.Logger.
func (l GenericSugar) Level() types.Level {
	return l.CompactLogger.Level()
}

// Log implements logger.Logger.
func (l GenericSugar) Log(level types.Level, values ...any) {
	l.CompactLogger.Log(level, values...)
}

// LogFields implements logger.Logger.
func (l GenericSugar) LogFields(level types.Level, message string, fields field.AbstractFields) {
	l.CompactLogger.LogFields(level, message, fields)
}

// Logf implements logger.Logger.
func (l GenericSugar) Logf(level types.Level, format string, args ...any) {
	l.CompactLogger.Logf(level, format, args...)
}

// Flush implements logger.Logger.
func (l GenericSugar) Flush() {
	l.CompactLogger.Flush()
}

// Emitter implements logger.Logger.
func (l GenericSugar) Emitter() types.Emitter {
	return l.CompactLogger.Emitter()
}

// WithLevel implements logger.Logger.
func (l GenericSugar) WithLevel(newLevel types.Level) types.Logger {
	return GenericSugar{CompactLogger: l.CompactLogger.WithLevel(newLevel)}
}

// WithPreHooks implements logger.Logger.
func (l GenericSugar) WithPreHooks(preHooks ...types.PreHook) types.Logger {
	return GenericSugar{CompactLogger: l.CompactLogger.WithPreHooks(preHooks...)}
}

// WithHooks implements logger.Logger.
func (l GenericSugar) WithHooks(hooks ...types.Hook) types.Logger {
	return GenericSugar{CompactLogger: l.CompactLogger.WithHooks(hooks...)}
}

// WithField implements logger.Logger.
func (l GenericSugar) WithField(
	key field.Key,
	value field.Value,
	props ...field.Property,
) types.Logger {
	return GenericSugar{CompactLogger: l.CompactLogger.WithField(key, value, props...)}
}

// WithFields implements logger.Logger.
func (l GenericSugar) WithFields(fields field.AbstractFields) types.Logger {
	return GenericSugar{CompactLogger: l.CompactLogger.WithFields(fields)}
}

// WithTraceIDs implements logger.Logger.
func (l GenericSugar) WithTraceIDs(allTraceIDs belt.TraceIDs, newTraceIDs int) belt.Tool {
	return GenericSugar{CompactLogger: l.CompactLogger.WithTraceIDs(allTraceIDs, newTraceIDs).(CompactLogger)}
}

// WithContextFields implements logger.Logger.
func (l GenericSugar) WithContextFields(allFields *field.FieldsChain, newFields int) belt.Tool {
	return GenericSugar{CompactLogger: l.CompactLogger.WithContextFields(allFields, newFields).(CompactLogger)}
}

// WithMessagePrefix implements logger.Logger.
func (l GenericSugar) WithMessagePrefix(prefix string) types.Logger {
	return GenericSugar{CompactLogger: l.CompactLogger.WithMessagePrefix(prefix)}
}

// TraceFields implements logger.Logger.
func (l GenericSugar) TraceFields(message string, fields field.AbstractFields) {
	l.LogFields(types.LevelTrace, message, fields)
}

// DebugFields implements logger.Logger.
func (l GenericSugar) DebugFields(message string, fields field.AbstractFields) {
	l.LogFields(types.LevelDebug, message, fields)
}

// InfoFields implements logger.Logger.
func (l GenericSugar) InfoFields(message string, fields field.AbstractFields) {
	l.LogFields(types.LevelInfo, message, fields)
}

// WarnFields implements logger.Logger.
func (l GenericSugar) WarnFields(message string, fields field.AbstractFields) {
	l.LogFields(types.LevelWarning, message, fields)
}

// ErrorFields implements logger.Logger.
func (l GenericSugar) ErrorFields(message string, fields field.AbstractFields) {
	l.LogFields(types.LevelError, message, fields)
}

// PanicFields implements logger.Logger.
func (l GenericSugar) PanicFields(message string, fields field.AbstractFields) {
	l.LogFields(types.LevelPanic, message, fields)
}

// FatalFields implements logger.Logger.
func (l GenericSugar) FatalFields(message string, fields field.AbstractFields) {
	l.LogFields(types.LevelFatal, message, fields)
}

// Trace implements logger.Logger.
func (l GenericSugar) Trace(values ...any) {
	l.Log(types.LevelTrace, values...)
}

// Debug implements logger.Logger.
func (l GenericSugar) Debug(values ...any) {
	l.Log(types.LevelDebug, values...)
}

// Info implements logger.Logger.
func (l GenericSugar) Info(values ...any) {
	l.Log(types.LevelInfo, values...)
}

// Warn implements logger.Logger.
func (l GenericSugar) Warn(values ...any) {
	l.Log(types.LevelWarning, values...)
}

// Error implements logger.Logger.
func (l GenericSugar) Error(values ...any) {
	l.Log(types.LevelError, values...)
}

// Panic implements logger.Logger.
func (l GenericSugar) Panic(values ...any) {
	l.Log(types.LevelPanic, values...)
}

// Fatal implements logger.Logger.
func (l GenericSugar) Fatal(values ...any) {
	l.Log(types.LevelFatal, values...)
}

// Tracef implements logger.Logger.
func (l GenericSugar) Tracef(format string, args ...any) {
	l.Logf(types.LevelTrace, format, args...)
}

// Debugf implements logger.Logger.
func (l GenericSugar) Debugf(format string, args ...any) {
	l.Logf(types.LevelDebug, format, args...)
}

// Infof implements logger.Logger.
func (l GenericSugar) Infof(format string, args ...any) {
	l.Logf(types.LevelInfo, format, args...)
}

// Warnf implements logger.Logger.
func (l GenericSugar) Warnf(format string, args ...any) {
	l.Logf(types.LevelWarning, format, args...)
}

// Errorf implements logger.Logger.
func (l GenericSugar) Errorf(format string, args ...any) {
	l.Logf(types.LevelError, format, args...)
}

// Panicf implements logger.Logger.
func (l GenericSugar) Panicf(format string, args ...any) {
	l.Logf(types.LevelPanic, format, args...)
}

// Fatalf implements logger.Logger.
func (l GenericSugar) Fatalf(format string, args ...any) {
	l.Logf(types.LevelFatal, format, args...)
}
