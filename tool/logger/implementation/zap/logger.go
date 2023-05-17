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

package zap

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/facebookincubator/go-belt"
	"github.com/facebookincubator/go-belt/pkg/field"
	"github.com/facebookincubator/go-belt/pkg/valuesparser"
	"github.com/facebookincubator/go-belt/tool/logger"
	"github.com/facebookincubator/go-belt/tool/logger/adapter"
	"github.com/facebookincubator/go-belt/tool/logger/types"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// FieldNameTraceIDs is the field name used to store belt.TraceIDs.
	FieldNameTraceIDs = "trace_id"
)

// Emitter is the implementation of types.Emitter based on a zap (non-sugared) logger.
type Emitter struct {
	ZapLogger *zap.Logger
}

var _ types.Emitter = (*Emitter)(nil)

// NewEmitter returns a new instance of Emitter
func NewEmitter(zapLogger *zap.Logger) Emitter {
	return Emitter{
		ZapLogger: zapLogger,
	}
}

const (
	preallocateFieldsSize = 16
)

var fieldsPool = sync.Pool{
	New: func() any {
		fields := make([]zap.Field, 0, preallocateFieldsSize)
		return &fields
	},
}

func acquireZapFields() *[]zap.Field {
	return fieldsPool.Get().(*[]zap.Field)
}

func releaseZapFields(fieldsPtr *[]zap.Field) {
	if fieldsPtr == nil || len(*fieldsPtr) >= 256 {
		return
	}
	*fieldsPtr = (*fieldsPtr)[:0]
	fieldsPool.Put(fieldsPtr)
}

// Flush implements types.Emitter.
func (l Emitter) Flush() {
	err := l.getZapLogger().Sync()
	if err == nil {
		return
	}

	errDesc := err.Error()
	if strings.Contains(errDesc, "sync /dev/stderr: invalid argument") {
		// Intentionally ignore this error, see https://github.com/uber-go/zap/issues/328
		return
	}

	l.Emit(&types.Entry{
		Timestamp: time.Now(),
		Level:     types.LevelError,
		Message:   fmt.Sprintf("unable to Sync(): %s", errDesc),
	})
}

// CheckZapLevel returns true if the given zap's logging level is enabled in this logger.
func (l Emitter) CheckZapLevel(level zapcore.Level) bool {
	return l.getZapLogger().Core().Enabled(level)
}

// CheckLevel returns true if the given types.Level is enabled in this logger.
func (l Emitter) CheckLevel(level types.Level) bool {
	return l.CheckZapLevel(LevelToZap(level))
}

// Emit implements types.Emitter.
func (l Emitter) Emit(entry *types.Entry) {
	// We log panics and fatals even if logging level is lower because
	// they supposed to trigger panic or os.Exit, and we cannot just return without that.
	// So it should be silent, but destructive.
	if entry.Level == types.LevelNone || !l.CheckLevel(entry.Level) && entry.Level > types.LevelPanic {
		return
	}

	var zapFields []zap.Field
	if !entry.Properties.Has(EntryPropertyIgnoreFields) && entry.Fields != nil {
		zapFieldsPtr := fieldsToZap(entry.Fields.Len(), entry.Fields.ForEachField)
		if zapFieldsPtr != nil {
			defer releaseZapFields(zapFieldsPtr)
			zapFields = *zapFieldsPtr
		}
	}

	zapEntry := zapcore.Entry{
		Level:      LevelToZap(entry.Level),
		Time:       entry.Timestamp,
		LoggerName: "",
		Message:    entry.Message,
		Stack:      "",
	}
	if entry.Caller != 0 {
		file, line := entry.Caller.FileLine()
		zapEntry.Caller = zapcore.EntryCaller{
			Defined: true,
			PC:      uintptr(entry.Caller),
			File:    file,
			Line:    line,
		}
	}

	l.LogZapEntry(zapEntry, zapFields...)
}

// LogZapEntry logs an entry given its zap structure.
func (l Emitter) LogZapEntry(zapEntry zapcore.Entry, zapFields ...zap.Field) {
	l.ZapLogger.Core().Check(zapEntry, nil).Write(zapFields...)
}

type mostlyPersistentData struct {
	entryPool     *sync.Pool
	zapEntryPool  *sync.Pool
	fmtBufPool    *sync.Pool
	preHooks      logger.PreHooks
	hooks         logger.Hooks
	traceIDs      belt.TraceIDs
	getCallerFunc types.GetCallerPC
	messagePrefix string
}

// CompactLogger is an implementation of types.CompactLogger based on a zap logger.
type CompactLogger struct {
	*mostlyPersistentData
	emitter            Emitter
	contextFields      *field.FieldsChain
	contextNewFields   uint32
	prepareEmitterOnce sync.Once
}

var _ adapter.CompactLogger = (*CompactLogger)(nil)

// Flush implements types.CompactLogger.
func (l *CompactLogger) Flush() {
	l.emitter.Flush()
	for _, hook := range l.hooks {
		hook.Flush()
	}
}

var timeNow = time.Now

var entryPropertiesIgnoreFields = types.EntryProperties{EntryPropertyIgnoreFields}

func (l *CompactLogger) acquireEntry() *types.Entry {
	entry := l.entryPool.Get().(*types.Entry)
	entry.Timestamp = timeNow()
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

func (l *CompactLogger) acquireZapEntry() *zapcore.Entry {
	entry := l.zapEntryPool.Get().(*zapcore.Entry)
	entry.Time = timeNow()
	return entry
}

func (l *CompactLogger) releaseZapEntry(entry *zapcore.Entry) {
	entry.Message = ""
	entry.Stack = ""
	entry.Caller.Defined = false
	l.zapEntryPool.Put(entry)
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

func (l *CompactLogger) logEntry(entry *logger.Entry) {
	if !entry.Caller.Defined() && l.getCallerFunc != nil {
		entry.Caller = l.getCallerFunc()
	}
	if !adapter.ProcessHooks(l.hooks, entry) {
		return
	}

	zapEntry := l.acquireZapEntry()
	defer l.releaseZapEntry(zapEntry)

	zapEntry.Time = entry.Timestamp
	zapEntry.Caller.Defined = entry.Caller.Defined()
	if zapEntry.Caller.Defined {
		file, line := entry.Caller.FileLine()
		zapEntry.Caller = zapcore.EntryCaller{
			Defined: true,
			PC:      uintptr(entry.Caller),
			File:    file,
			Line:    line,
		}
	}
	zapEntry.Message = entry.Message
	zapEntry.Level = LevelToZap(entry.Level)
	var zapFields []zap.Field
	zapFieldsPtr := fieldsToZap(entry.Fields.Len(), entry.Fields.ForEachField)
	if zapFieldsPtr != nil {
		defer releaseZapFields(zapFieldsPtr)
		zapFields = *zapFieldsPtr
	}

	l.logZapEntryNoHooks(zapEntry, zapFields...)
}

func (l *CompactLogger) logZapEntryNoHooks(zapEntry *zapcore.Entry, zapFields ...zapcore.Field) {
	l.prepareEmitter()
	l.zapSetCaller(zapEntry)
	l.emitter.LogZapEntry(*zapEntry, zapFields...)

	// zap authors decided not to panic or/and exit on Panics and Fatals,
	// see: https://github.com/uber-go/zap/issues/358
	//
	// So doing this manually here:
	switch zapEntry.Level {
	case zap.PanicLevel, zap.DPanicLevel:
		panic(fmt.Sprintf("%#+v: %#+v", *zapEntry, zapFields))
	case zap.FatalLevel:
		os.Exit(1)
	}
}

func (l *CompactLogger) zapSetCaller(zapEntry *zapcore.Entry) {
	if zapEntry.Caller.Defined || l.getCallerFunc == nil {
		return
	}

	caller := l.getCallerFunc()
	if !caller.Defined() {
		return
	}
	file, line := caller.FileLine()

	zapEntry.Caller = zapcore.NewEntryCaller(
		uintptr(caller),
		file,
		line,
		true,
	)
}

// LogFields implements types.CompactLogger.
func (l *CompactLogger) LogFields(level types.Level, message string, fields field.AbstractFields) {
	preHooksResult := adapter.LogFieldsPreprocess(l.preHooks, l.traceIDs, level, l.emitter.CheckLevel(level), message, fields)
	if preHooksResult.Skip {
		return
	}

	if preHooksResult.ExtraFields != nil {
		fields = field.Slice[field.AbstractFields]{fields, preHooksResult.ExtraFields}
	}

	if len(l.hooks) != 0 {
		entry := l.acquireEntry()
		defer l.releaseEntry(entry)

		entry.Level = level
		entry.Message = message
		entry.Fields = fields
		l.logEntry(entry)
		return
	}

	entry := l.acquireZapEntry()
	defer l.releaseZapEntry(entry)

	entry.Level = LevelToZap(level)
	entry.Message = message

	var (
		zapFields    []zap.Field
		zapFieldsPtr *[]zap.Field
	)
	zapFieldsPtr = fieldsToZap(fields.Len(), fields.ForEachField)
	if zapFieldsPtr != nil {
		defer releaseZapFields(zapFieldsPtr)
		zapFields = *zapFieldsPtr
	}
	l.logZapEntryNoHooks(entry, zapFields...)
}

// Logf implements types.CompactLogger.
func (l *CompactLogger) Logf(level types.Level, format string, args ...any) {
	preHooksResult := adapter.LogfPreprocess(l.preHooks, l.traceIDs, level, l.emitter.CheckLevel(level), format, args...)
	if preHooksResult.Skip {
		return
	}

	buf := l.acquireBuf()
	defer l.releaseBuf(buf)
	fmt.Fprintf(buf, format, args...)

	if len(l.hooks) != 0 {
		entry := l.acquireEntry()
		defer l.releaseEntry(entry)

		entry.Level = level
		entry.Message = buf.String()
		entry.Fields = preHooksResult.ExtraFields
		l.logEntry(entry)
		return
	}

	entry := l.acquireZapEntry()
	defer l.releaseZapEntry(entry)

	entry.Level = LevelToZap(level)
	entry.Message = buf.String()

	var zapFields []zap.Field
	if preHooksResult.ExtraFields != nil {
		zapFieldsPtr := fieldsToZap(preHooksResult.ExtraFields.Len(), preHooksResult.ExtraFields.ForEachField)
		if zapFieldsPtr != nil {
			defer releaseZapFields(zapFieldsPtr)
			zapFields = *zapFieldsPtr
		}
	}

	l.logZapEntryNoHooks(entry, zapFields...)
}

// Log implements types.CompactLogger.
func (l *CompactLogger) Log(level types.Level, values ...any) {
	forceProcess := level == logger.LevelFatal || level == logger.LevelPanic

	preHooksResult := adapter.LogPreprocess(l.preHooks, l.traceIDs, level, l.emitter.CheckLevel(level), values...)
	if preHooksResult.Skip && !forceProcess {
		return
	}

	if len(values) == 1 {
		if entry, ok := values[0].(*logger.Entry); ok {
			if entry.TraceIDs != nil {
				entry.Fields = field.Add(entry.Fields, preHooksResult.ExtraFields, &field.Field{Key: FieldNameTraceIDs, Value: entry.TraceIDs})
			} else {
				entry.Fields = field.Add(entry.Fields, preHooksResult.ExtraFields)
			}
			l.logEntry(entry)
			return
		}
	}

	valuesParser := valuesparser.AnySlice(values)

	if len(l.hooks) != 0 {
		entry := l.acquireEntry()
		defer l.releaseEntry(entry)

		entry.Level = level
		entry.Fields = field.Add(entry.Fields, preHooksResult.ExtraFields)

		buf := l.acquireBuf()
		defer l.releaseBuf(buf)
		valuesParser.ForEachField(func(f *field.Field) bool { return true }) // just causing to parse the fields to make valuesParser.WriteUnparsed(buf) work correctly
		entry.Message = buf.String()

		l.logEntry(entry)
		return
	}

	entry := l.acquireZapEntry()
	defer l.releaseZapEntry(entry)

	entry.Level = LevelToZap(level)

	var (
		zapFields    []zap.Field
		zapFieldsPtr *[]zap.Field
	)
	if preHooksResult.ExtraFields != nil {
		fields := field.Slice[field.AbstractFields]{&valuesParser, preHooksResult.ExtraFields}
		zapFieldsPtr = fieldsToZap(fields.Len(), fields.ForEachField)
	} else {
		zapFieldsPtr = fieldsToZap(valuesParser.Len(), valuesParser.ForEachField)
	}
	if zapFieldsPtr != nil {
		defer releaseZapFields(zapFieldsPtr)
		zapFields = *zapFieldsPtr
	}

	buf := l.acquireBuf()
	defer l.releaseBuf(buf)
	valuesParser.WriteUnparsed(buf)
	entry.Message = buf.String()

	l.logZapEntryNoHooks(entry, zapFields...)
}

// Level implements types.CompactLogger.
func (l *CompactLogger) Level() types.Level {
	core := l.emitter.ZapLogger.Core()
	for level := types.EndOfLevel - 1; level >= 0; level-- {
		if core.Enabled(LevelToZap(level)) {
			return level
		}
	}
	return types.LevelNone
}

// WithLevel implements types.CompactLogger.
func (l *CompactLogger) WithLevel(newLevel types.Level) adapter.CompactLogger {
	branch := l.branch()
	branch.emitter.ZapLogger = branch.emitter.ZapLogger.WithOptions(zap.IncreaseLevel(LevelToZap(newLevel)))
	return branch
}

func (l *CompactLogger) branch() *CompactLogger {
	return &CompactLogger{
		mostlyPersistentData: l.mostlyPersistentData,
		emitter:              Emitter{ZapLogger: l.emitter.getZapLogger()},
		contextFields:        l.contextFields,
		contextNewFields:     atomic.LoadUint32(&l.contextNewFields),
	}
}

func (l *CompactLogger) clone() *CompactLogger {
	clone := l.branch()
	clone.mostlyPersistentData = &[]mostlyPersistentData{*l.mostlyPersistentData}[0]
	return clone
}

// WithMessagePrefix implements types.CompactLogger.
func (l *CompactLogger) WithMessagePrefix(prefix string) adapter.CompactLogger {
	clone := l.clone()
	clone.messagePrefix += prefix
	return clone
}

// WithHooks implements types.CompactLogger.
func (l *CompactLogger) WithHooks(hooks ...types.Hook) adapter.CompactLogger {
	clone := l.clone()
	clone.hooks = make(types.Hooks, len(l.hooks)+len(hooks))
	copy(clone.hooks, l.hooks)
	copy(clone.hooks[len(l.hooks):], hooks)
	return clone
}

// WithPreHooks implements types.CompactLogger.
func (l *CompactLogger) WithPreHooks(preHooks ...types.PreHook) adapter.CompactLogger {
	clone := l.clone()
	clone.preHooks = make(types.PreHooks, len(l.preHooks)+len(preHooks))
	copy(clone.preHooks, l.preHooks)
	copy(clone.preHooks[len(l.preHooks):], preHooks)
	return clone
}

// WithField implements types.CompactLogger.
func (l *CompactLogger) WithField(
	key field.Key,
	value field.Value,
	props ...field.Property,
) adapter.CompactLogger {
	branch := l.branch()
	branch.contextFields = l.contextFields.WithField(key, value, props)
	branch.contextNewFields++
	return branch
}

// WithFields implements types.CompactLogger.
func (l *CompactLogger) WithFields(fields field.AbstractFields) adapter.CompactLogger {
	branch := l.branch()
	branch.contextFields = l.contextFields.WithFields(fields)
	branch.contextNewFields += uint32(fields.Len())
	return branch
}

// WithTraceIDs implements types.CompactLogger and belt.Tool.
func (l *CompactLogger) WithTraceIDs(allTraceIDs belt.TraceIDs, newTraceIDsCount int) belt.Tool {
	clone := l.clone()
	clone.traceIDs = allTraceIDs
	clone.contextFields = clone.contextFields.WithField(FieldNameTraceIDs, allTraceIDs)
	clone.contextNewFields++
	return clone
}

// WithContextFields implements types.CompactLogger and belt.Tool.
func (l *CompactLogger) WithContextFields(allFields *field.FieldsChain, newFieldsCount int) belt.Tool {
	branch := l.branch()
	branch.contextFields = allFields
	branch.contextNewFields += uint32(newFieldsCount)
	return branch
}

func (l *CompactLogger) compileEmitterFields() {
	// This function is the only one which can change
	// the value of l.contextNewFields, and it is called
	// only within l.prepareEmitterOnce.Do, so:
	// * We do not need to atomically load the l.contextNewFields value.
	// * We still need to atomically store the l.contextNewFields value,
	//   because other goroutines may read it concurrently in other functions.
	if l.contextNewFields == 0 {
		return
	}
	defer atomic.StoreUint32(&l.contextNewFields, 0)
	zapFieldsPtr := fieldsToZap(int(l.contextNewFields), l.contextFields.ForEachField)
	if zapFieldsPtr == nil {
		return
	}
	defer releaseZapFields(zapFieldsPtr)
	l.emitter.setZapLogger(l.emitter.ZapLogger.With(*zapFieldsPtr...))
}

func (l *CompactLogger) prepareEmitter() {
	l.prepareEmitterOnce.Do(func() {
		l.compileEmitterFields()
	})
}

// Emitter implements types.CompactLogger.
func (l *CompactLogger) Emitter() types.Emitter {
	l.prepareEmitter()
	return l.emitter
}

// DefaultZapLogger is the (overridable) function which returns
// a zap logger with the default configuration.
//
// Do not override this anywhere but in the `main` package.
var DefaultZapLogger = func() *zap.Logger {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(LevelToZap(logger.LevelTrace))
	zapLogger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return zapLogger
}

// Default returns a logger.Logger using the default zap logger
// (see DefaultZapLogger).
func Default() types.Logger {
	return New(DefaultZapLogger())
}

// New returns a new instance of logger.Logger given a zap logger.
func New(logger *zap.Logger, opts ...types.Option) types.Logger {
	return adapter.GenericSugar{
		CompactLogger: newCompactLoggerFromZap(logger, opts...),
	}
}

func newCompactLoggerFromZap(zapLogger *zap.Logger, opts ...types.Option) *CompactLogger {
	cfg := types.Options(opts).Config()
	return &CompactLogger{
		emitter: NewEmitter(zapLogger),
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
			zapEntryPool: &sync.Pool{
				New: func() any {
					return &zapcore.Entry{}
				},
			},
		},
	}
}
