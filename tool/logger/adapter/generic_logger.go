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
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/facebookincubator/go-belt"
	"github.com/facebookincubator/go-belt/pkg/field"
	"github.com/facebookincubator/go-belt/tool/logger/experimental"
	"github.com/facebookincubator/go-belt/tool/logger/types"
)

// GenericLogger is a generic implementation of types.CompactLogger given
// Emitters.
type GenericLogger struct {
	Emitters        types.Emitters
	CurrentLevel    types.Level
	Fields          *field.FieldsChain
	TraceIDs        belt.TraceIDs
	CurrentPreHooks types.PreHooks
	CurrentHooks    types.Hooks
	MessagePrefix   string
	GetCallerFunc   types.GetCallerPC
}

var _ CompactLogger = (*GenericLogger)(nil)

// Level implements types.CompactLogger.
func (l *GenericLogger) Level() types.Level {
	return l.CurrentLevel
}

// WithTraceIDs implements types.CompactLogger.
func (l *GenericLogger) WithTraceIDs(allTraceIDs belt.TraceIDs, newTraceIDs int) belt.Tool {
	newLogger := *l
	newLogger.TraceIDs = allTraceIDs
	return &newLogger
}

// WithLevel implements types.CompactLogger.
func (l *GenericLogger) WithLevel(newLevel types.Level) CompactLogger {
	newLogger := *l
	newLogger.CurrentLevel = newLevel
	return &newLogger
}

// WithPreHooks implements types.CompactLogger.
func (l *GenericLogger) WithPreHooks(preHooks ...types.PreHook) CompactLogger {
	newLogger := *l
	if preHooks == nil {
		newLogger.CurrentPreHooks = nil
	} else {
		newLogger.CurrentPreHooks = types.PreHooks{newLogger.CurrentPreHooks, types.PreHooks(preHooks)}
	}
	return &newLogger
}

// WithHooks implements types.CompactLogger.
func (l *GenericLogger) WithHooks(hooks ...types.Hook) CompactLogger {
	newLogger := *l
	if hooks == nil {
		newLogger.CurrentHooks = nil
	} else {
		newLogger.CurrentHooks = types.Hooks{newLogger.CurrentHooks, types.Hooks(hooks)}
	}
	return &newLogger
}

// WithField implements types.CompactLogger.
func (l *GenericLogger) WithField(
	key field.Key,
	value field.Value,
	props ...field.Property,
) CompactLogger {
	newLogger := *l
	newLogger.Fields = newLogger.Fields.WithField(key, value, props...)
	return &newLogger
}

// WithFields implements types.CompactLogger.
func (l *GenericLogger) WithFields(fields field.AbstractFields) CompactLogger {
	newLogger := *l
	if fields == nil {
		return &newLogger
	}
	newLogger.Fields = newLogger.Fields.WithFields(fields)
	return &newLogger
}

// WithContextFields implements types.CompactLogger.
func (l *GenericLogger) WithContextFields(allFields *field.FieldsChain, newFields int) belt.Tool {
	newLogger := *l
	newLogger.Fields = allFields
	return &newLogger
}

// Emitter implements types.CompactLogger.
func (l *GenericLogger) Emitter() types.Emitter {
	return l.Emitters
}

// ValuesParser a handler for arbitrary values, which extracts structured fields/values
// by method ForEachField() and provides everything else as a string by method WriteUnparsed.
//
// It also implements fields.AbstractFields, so it could be directly used as a collection
// of structured fields right after just a type-casting ([]any -> ValuesParser). It is
// supposed to be a zero allocation implementation in future.
//
// For example it is used to parse all the arguments of logger.Logger.Log.
type ValuesParser []any

// ForEachField implements field.AbstractFields.
func (p *ValuesParser) ForEachField(callback func(f *field.Field) bool) bool {
	for idx := 0; idx < len(*p); {
		value := (*p)[idx]

		if value == nil {
			idx++
			continue
		}
		v := reflect.Indirect(reflect.ValueOf(value))
		if forEachFielder, ok := v.Interface().(field.ForEachFieldser); ok {
			if !forEachFielder.ForEachField(callback) {
				return false
			}
			idx++
			continue
		}

		switch v.Kind() {
		case reflect.Map:
			r := ParseMap(v, callback)
			(*p)[idx] = nil
			if !r {
				return false
			}
		case reflect.Struct:
			r := ParseStruct(nil, v, callback)
			(*p)[idx] = nil
			if !r {
				return false
			}
		default:
			idx++
			continue
		}
	}
	return true
}

// ParseMap calls the callback for each pair in the map until first false is returned
//
// It returns false if callback returned false.
func ParseMap(m reflect.Value, callback func(f *field.Field) bool) bool {
	var f field.Field
	for _, keyV := range m.MapKeys() {
		valueV := m.MapIndex(keyV)
		switch key := keyV.Interface().(type) {
		case field.Key:
			f.Key = key
		default:
			f.Key = fmt.Sprint(key)
		}
		f.Value = valueV.Interface()
		if !callback(&f) {
			return false
		}
	}
	return true
}

// ParseStruct parses a structure to a collection of fields.
//
// `fieldPath` is the prefix of the field-name.
// `_struct` is the structure to be parsed (provided as a reflect.Value).
// `callback` is the function called for each found field, until first false is returned.
//
// It returns false if callback returned false.
func ParseStruct(fieldPath []string, _struct reflect.Value, callback func(f *field.Field) bool) bool {
	s := reflect.Indirect(_struct)

	var f field.Field
	// TODO: optimize this
	t := s.Type()

	fieldPath = append(fieldPath, "")

	fieldCount := s.NumField()
	for fieldNum := 0; fieldNum < fieldCount; fieldNum++ {
		structFieldType := t.Field(fieldNum)
		if structFieldType.PkgPath != "" {
			// unexported
			continue
		}
		logTag := structFieldType.Tag.Get("log")
		if logTag == "-" {
			continue
		}
		structField := s.Field(fieldNum)
		if structField.IsZero() {
			continue
		}
		value := reflect.Indirect(structField)

		pathComponent := structFieldType.Name
		if logTag != "" {
			pathComponent = logTag
		}
		fieldPath[len(fieldPath)-1] = pathComponent

		if value.Kind() == reflect.Struct {
			if !ParseStruct(fieldPath, value, callback) {
				return false
			}
			continue
		}

		f.Key = strings.Join(fieldPath, ".")
		f.Value = value.Interface()
		if !callback(&f) {
			return false
		}
	}

	return true
}

// Len implements field.AbstractFields.
func (p *ValuesParser) Len() int {
	return len(*p)
}

// WriteUnparsed writes unstructud data (everything that was never considered as a structued field by
// ForEachField method) to the given io.Writer.
func (p *ValuesParser) WriteUnparsed(w io.Writer) {
	for _, value := range *p {
		if value == nil {
			continue
		}

		fmt.Fprint(w, value)
	}
}

// Log implements types.CompactLogger.
func (l *GenericLogger) Log(level types.Level, values ...any) {
	preHooksResult := LogPreprocess(l.CurrentPreHooks, l.TraceIDs, level, l.CurrentLevel >= level, values...)
	if preHooksResult.Skip {
		return
	}

	if len(values) == 1 {
		if entry, ok := values[0].(*types.Entry); ok {
			entry.Fields = field.Add(entry.Fields, preHooksResult.ExtraFields)
			l.emit(entry)
			return
		}
	}

	// TODO: optimize this
	valuesParser := ValuesParser(values)
	fields := make(field.Fields, 0, valuesParser.Len())
	valuesParser.ForEachField(func(f *field.Field) bool {
		fields = append(fields, *f)
		return true
	})

	// TODO: optimize this
	var buf strings.Builder
	valuesParser.WriteUnparsed(&buf)

	var finalFields field.AbstractFields
	if preHooksResult.ExtraFields != nil {
		finalFields = field.Slice[field.AbstractFields]{fields, preHooksResult.ExtraFields}
	} else {
		finalFields = fields
	}

	entry := l.acquireEntry(level, buf.String(), finalFields, preHooksResult.ExtraEntryProperties)
	defer releaseEntry(entry)

	l.emit(entry)
}

// LogFields implements types.CompactLogger.
func (l *GenericLogger) LogFields(level types.Level, message string, fields field.AbstractFields) {
	preHooksResult := LogFieldsPreprocess(l.CurrentPreHooks, l.TraceIDs, level, l.CurrentLevel >= level, message, fields)
	if preHooksResult.Skip {
		return
	}

	var finalFields field.AbstractFields
	if preHooksResult.ExtraFields != nil {
		finalFields = field.Slice[field.AbstractFields]{fields, preHooksResult.ExtraFields}
	} else {
		finalFields = fields
	}

	entry := l.acquireEntry(level, message, finalFields, preHooksResult.ExtraEntryProperties)
	defer releaseEntry(entry)

	l.emit(entry)
}

// LogPreprocess checks logging level and calls PreHooks. Returns the total
// outcome of these two steps.
//
// It will never return Skip=true for message logging levels Panic and Fatal, instead
// it will ask to skip the logging through Entry property "experimental.EntryPropertySkipAllEmitters".
//
// It is a helper which is supposed to be used in the beginning of a Log implementation.
func LogPreprocess(
	preHooks types.PreHooks,
	traceIDs belt.TraceIDs,
	level types.Level,
	logLevelSatisfied bool,
	values ...any,
) types.PreHookResult {
	if level == types.LevelNone {
		return types.PreHookResult{Skip: true}
	}
	shouldSkip := false
	if !logLevelSatisfied {
		if level != types.LevelPanic && level != types.LevelFatal {
			return types.PreHookResult{Skip: true}
		}
		shouldSkip = true
	}

	result := preHooks.ProcessInput(traceIDs, level, values...)
	if !result.Skip {
		if shouldSkip {
			result.ExtraEntryProperties = append(result.ExtraEntryProperties, experimental.EntryPropertySkipAllEmitters)
		}
		return result
	}
	if level != types.LevelPanic && level != types.LevelFatal {
		return result
	}

	result.Skip = false
	result.ExtraEntryProperties = append(result.ExtraEntryProperties, experimental.EntryPropertySkipAllEmitters)
	return result
}

// LogFieldsPreprocess checks logging level and calls PreHooks. Returns the total
// outcome of these two steps.
//
// It will never return Skip=true for message logging levels Panic and Fatal, instead
// it will ask to skip the logging through Entry property "experimental.EntryPropertySkipAllEmitters".
//
// It is a helper which is supposed to be used in the beginning of a LogFields implementation.
func LogFieldsPreprocess(
	preHooks types.PreHooks,
	traceIDs belt.TraceIDs,
	level types.Level,
	logLevelSatisfied bool,
	message string,
	fields field.AbstractFields,
) types.PreHookResult {
	if level == types.LevelNone {
		return types.PreHookResult{Skip: true}
	}
	shouldSkip := false
	if !logLevelSatisfied {
		if level != types.LevelPanic && level != types.LevelFatal {
			return types.PreHookResult{Skip: true}
		}
		shouldSkip = true
	}

	result := preHooks.ProcessInputFields(traceIDs, level, message, fields)
	if !result.Skip {
		if shouldSkip {
			result.ExtraEntryProperties = append(result.ExtraEntryProperties, experimental.EntryPropertySkipAllEmitters)
		}
		return result
	}
	if level != types.LevelPanic && level != types.LevelFatal {
		return result
	}

	result.Skip = false
	result.ExtraEntryProperties = append(result.ExtraEntryProperties, experimental.EntryPropertySkipAllEmitters)
	return result
}

// LogfPreprocess checks logging level and calls PreHooks. Returns the total
// outcome of these two steps.
//
// It will never return Skip=true for message with logging levels Panic and Fatal, instead
// it will ask to skip the logging through Entry property "experimental.EntryPropertySkipAllEmitters".
//
// It is a helper which is supposed to be used in the beginning of a Logf implementation.
func LogfPreprocess(
	preHooks types.PreHooks,
	traceIDs belt.TraceIDs,
	level types.Level,
	logLevelSatisfied bool,
	format string,
	args ...any,
) types.PreHookResult {
	if level == types.LevelNone {
		return types.PreHookResult{Skip: true}
	}
	shouldSkip := false
	if !logLevelSatisfied {
		if level != types.LevelPanic && level != types.LevelFatal {
			return types.PreHookResult{Skip: true}
		}
		shouldSkip = true
	}

	result := preHooks.ProcessInputf(traceIDs, level, format, args...)
	if !result.Skip {
		if shouldSkip {
			result.ExtraEntryProperties = append(result.ExtraEntryProperties, experimental.EntryPropertySkipAllEmitters)
		}
		return result
	}
	if level != types.LevelPanic && level != types.LevelFatal {
		return result
	}

	result.Skip = false
	result.ExtraEntryProperties = append(result.ExtraEntryProperties, experimental.EntryPropertySkipAllEmitters)
	return result
}

// Logf implements types.CompactLogger.
func (l *GenericLogger) Logf(level types.Level, format string, args ...any) {
	preHooksResult := LogfPreprocess(l.CurrentPreHooks, l.TraceIDs, level, l.CurrentLevel >= level, format, args...)
	if preHooksResult.Skip {
		return
	}

	entry := l.acquireEntry(level, fmt.Sprintf(format, args...), preHooksResult.ExtraFields, preHooksResult.ExtraEntryProperties)
	defer releaseEntry(entry)

	l.emit(entry)
}

// WithMessagePrefix implements types.CompactLogger.
func (l *GenericLogger) WithMessagePrefix(prefix string) CompactLogger {
	newLogger := *l
	newLogger.MessagePrefix = newLogger.MessagePrefix + prefix
	return &newLogger
}

// ProcessHooks executes hooks and never returns false for Panic and Fatal logging
// levels. In case of a "false" result from hooks for Panic or Fatal it adds
// EntryProperty "experimental.EntryPropertySkipAllEmitters" and returns true.
//
// It is a helper which is supposed to be used in Logf, LogFields and Log implementations.
func ProcessHooks(hooks types.Hooks, entry *types.Entry) bool {
	if hooks.ProcessLogEntry(entry) {
		return true
	}

	if entry.Level != types.LevelPanic && entry.Level != types.LevelFatal {
		return false
	}

	entry.Properties = append(entry.Properties, experimental.EntryPropertySkipAllEmitters)
	return true
}

func (l *GenericLogger) emit(entry *types.Entry) {
	if !ProcessHooks(l.CurrentHooks, entry) {
		return
	}

	if !entry.Properties.Has(experimental.EntryPropertySkipAllEmitters) {
		l.Emitters.Emit(entry)
	}

	switch entry.Level {
	case types.LevelPanic, types.LevelFatal:
		l.Flush()
		switch entry.Level {
		case types.LevelPanic:
			panic(fmt.Sprintf("panic was requested with the log entry: %#v", entry))
		case types.LevelFatal:
			os.Exit(2)
		}
	}
}

// Flush implements types.CompactLogger.
func (l *GenericLogger) Flush() {
	l.Emitters.Flush()
}

func (l *GenericLogger) acquireEntry(level types.Level, message string, fields field.AbstractFields, props types.EntryProperties) *types.Entry {
	entry := acquireEntry()
	entry.Level = level
	entry.Timestamp = time.Now()
	entry.TraceIDs = l.TraceIDs
	entry.Message = l.MessagePrefix + message
	entry.Fields = l.Fields.WithFields(fields)
	entry.Properties = props
	entry.Caller = l.GetCallerFunc()
	return entry
}
