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

// Package types of logger unifies different types of loggers into
// interfaces Logger. For example it allows to upgrade simple fmt.Printf
// to be a fully functional Logger. Therefore multiple wrappers are implemented
// here to provide different functions which could be missing in some loggers.
package types

import (
	"github.com/facebookincubator/go-belt"
	"github.com/facebookincubator/go-belt/pkg/field"
)

// Logger is an abstract generic structured logger. It is supposed
// to be one of the belt.Tools provided by the Belt. And therefore interface
// supposed to be used by end users of the `logger` package, so it already
// contains syntax sugar.
//
// All methods are thread-safe.
//
// For an optimized interface with the same functionality use interface CompactLogger.
type Logger interface {
	belt.Tool

	// Logf logs an unstructured message. Though, of course, all
	// contextual structured fields will also be logged.
	//
	// This method exists mostly for convenience, for people who
	// has not got used to proper structured logging, yet.
	// See `LogFields` and `Log`. If one have variables they want to
	// log, it is better for scalable observability to log them
	// as structured values, instead of injecting them into a
	// non-structured string.
	Logf(level Level, format string, args ...any)

	// LogFields logs structured fields with a explanation message.
	//
	// Anything that implements field.AbstractFields might be used
	// as a collection of fields to be logged.
	//
	// Examples:
	//
	// 	l.LogFields(logger.LevelDebug, "new_request", field.Fields{{Key: "user_id", Value: userID}, {Key: "group_id", Value: groupID}})
	// 	l.LogFields(logger.LevelInfo, "affected entries", field.Field{Key: "mysql_affected", Value: affectedRows})
	// 	l.LogFields(logger.LevelError, "unable to fetch user info", request) // where `request` implements field.AbstractFields
	//
	// Sometimes it is inconvenient to manually describe each field,
	// and for such cases see method `Log`.
	LogFields(level Level, message string, fields field.AbstractFields)

	// Log extracts structured fields from provided values, joins
	// the rest into an unstructured message and logs the result.
	//
	// This function provides convenience (relatively to LogFields)
	// at cost of a bit of performance.
	//
	// There are few ways to extract structured fields, which are
	// applied for each value from `values` (in descending priority order):
	// 1. If a `value` is an `*Entry` then the Entry is used (with its fields).
	//    This works only if this is the only argument. Otherwise it is
	//    threated as a simple structure (see point #3).
	// 2. If a `value` implements field.AbstractFields then ForEachField method
	//    is used (so it is become similar to LogFields).
	// 3. If a `value` is a structure (or a pointer to a structure) then
	//    fields of the structure are interpreted as structured fields
	//    to be logged (see explanation below).
	// 4. If a `value` is a map then fields a constructed out of this map.
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
	// 	user, err := getUser()
	// 	if err != nil {
	// 		l.Log(logger.LevelError, err)
	// 		return err
	// 	}
	// 	l.Log(logger.LevelDebug, "current user", user) // fields of structure "user" will be logged
	// 	l.Log(logger.LevelDebug, map[string]any{"user_id": user.ID, "group_id", user.GroupID})
	// 	l.Log(logger.LevelDebug, field.Fields{{Key: "user_id", Value: user.ID}, {Key: "group_id", Value: user.GroupID}})
	// 	l.Log(logger.LevelDebug, "current user ID is ", user.ID, " and group ID is ", user.GroupID) // will result into message "current user ID is 1234 and group ID is 5678".
	Log(level Level, values ...any)

	// Emitter returns the Emitter (see the description of interface "Emitter").
	Emitter() Emitter

	// Level returns the current logging level (see description of "Level").
	Level() Level

	// TraceFields is just a shorthand for LogFields(logger.LevelTrace, ...)
	TraceFields(message string, fields field.AbstractFields)

	// DebugFields is just a shorthand for LogFields(logger.LevelDebug, ...)
	DebugFields(message string, fields field.AbstractFields)

	// InfoFields is just a shorthand for LogFields(logger.LevelInfo, ...)
	InfoFields(message string, fields field.AbstractFields)

	// WarnFields is just a shorthand for LogFields(logger.LevelWarn, ...)
	WarnFields(message string, fields field.AbstractFields)

	// ErrorFields is just a shorthand for LogFields(logger.LevelError, ...)
	ErrorFields(message string, fields field.AbstractFields)

	// PanicFields is just a shorthand for LogFields(logger.LevelPanic, ...)
	//
	// Be aware: Panic level also triggers a `panic`.
	PanicFields(message string, fields field.AbstractFields)

	// FatalFields is just a shorthand for LogFields(logger.LevelFatal, ...)
	//
	// Be aware: Panic level also triggers an `os.Exit`.
	FatalFields(message string, fields field.AbstractFields)

	// Trace is just a shorthand for Log(logger.LevelTrace, ...)
	Trace(values ...any)

	// Debug is just a shorthand for Log(logger.LevelDebug, ...)
	Debug(values ...any)

	// Info is just a shorthand for Log(logger.LevelInfo, ...)
	Info(values ...any)

	// Warn is just a shorthand for Log(logger.LevelWarn, ...)
	Warn(values ...any)

	// Error is just a shorthand for Log(logger.LevelError, ...)
	Error(values ...any)

	// Panic is just a shorthand for Log(logger.LevelPanic, ...)
	//
	// Be aware: Panic level also triggers a `panic`.
	Panic(values ...any)

	// Fatal is just a shorthand for Log(logger.LevelFatal, ...)
	//
	// Be aware: Fatal level also triggers an `os.Exit`.
	Fatal(values ...any)

	// Tracef is just a shorthand for Logf(logger.LevelTrace, ...)
	Tracef(format string, args ...any)

	// Debugf is just a shorthand for Logf(logger.LevelDebug, ...)
	Debugf(format string, args ...any)

	// Infof is just a shorthand for Logf(logger.LevelInfo, ...)
	Infof(format string, args ...any)

	// Warnf is just a shorthand for Logf(logger.LevelWarn, ...)
	Warnf(format string, args ...any)

	// Errorf is just a shorthand for Logf(logger.LevelError, ...)
	Errorf(format string, args ...any)

	// Panicf is just a shorthand for Logf(logger.LevelPanic, ...)
	//
	// Be aware: Panic level also triggers a `panic`.
	Panicf(format string, args ...any)

	// Fatalf is just a shorthand for Logf(logger.LevelFatal, ...)
	//
	// Be aware: Fatal level also triggers an `os.Exit`.
	Fatalf(format string, args ...any)

	// WithLevel returns a logger with logger level set to the given argument.
	//
	// See also the description of type "Level".
	WithLevel(Level) Logger

	// WithPreHooks returns a Logger which includes/appends pre-hooks from the arguments.
	//
	// See also description of "PreHook".
	//
	// Special case: to reset hooks use `WithPreHooks()` (without any arguments).
	WithPreHooks(...PreHook) Logger

	// WithHooks returns a Logger which includes/appends hooks from the arguments.
	//
	// See also description of "Hook".
	//
	// Special case: to reset hooks use `WithHooks()` (without any arguments).
	WithHooks(...Hook) Logger

	// WithField returns the logger with the added field (used for structured logging).
	WithField(key string, value any, props ...field.Property) Logger

	// WithFields returns the logger with the added fields (used for structured logging)
	WithFields(fields field.AbstractFields) Logger

	// WithMessagePrefix adds a string to all messages logged through the derived logger.
	WithMessagePrefix(prefix string) Logger

	// WithEntryProperties adds props to EntryProperties of each emitted Entry.
	// This could be used only for enabling implementation-specific behavior.
	WithEntryProperties(props ...EntryProperty) Logger
}

// Emitter is a log entry sender. It is not obligated to provide
// functionality for logging levels, hooks, contextual fields or
// any other fancy stuff, and it just sends what is was told to.
//
// Note:
// Some specific Emitter implementations may support filtering of
// messages based on log level or/and add structured fields
// internally or do other stuff. But it is expected that even
// if a Emitter actually supports that kind of functionality,
// it will still be by default configured in a way like
// it has no such functionality (thus maximal logging level, no contextual
// fields and so on).
// However, if a Emitter is returned from a Logger, then it may (or may not)
// inherit properties of Logger (like logging level or structured fields).
// The undefined behavior here is left intentionally to provide more flexibility
// to Logger implementations to achieve better performance. It is considered
// than any Emitter managed by a Logger may have any configuration
// the Logger may consider the optimal one at any moment.
//
// All methods are thread-safe.
type Emitter interface {
	Flusher

	// Emit just logs the provided Entry. It does not modify it.
	//
	// If it is reasonably possible then the implementation of Emitter
	// should not panic or os.Exit even if the Level is Fatal or Panic.
	// Otherwise for example it prevents from logging a problem through other
	// Emitters if there are multiple of them.
	Emit(entry *Entry)
}

// Emitters is a set of Emitter-s.
//
// Only the last Emitter is allowed to panic or/and os.Exit (on Level-s Fatal and Panic).
// Do not use a set of Emitters if this rule is not satisfied. Check guarantees these
// provided by specific a Emitter implementation.
type Emitters []Emitter

var _ Emitter = (Emitters)(nil)

// Flush implements Emitter.
func (s Emitters) Flush() {
	for _, l := range s {
		l.Flush()
	}
}

// Emit implements Emitter.
func (s Emitters) Emit(entry *Entry) {
	for _, l := range s {
		l.Emit(entry)
	}
}

// Flusher defines a method to flush all buffers.
type Flusher interface {
	// Flush forces to flush all buffers.
	Flush()
}
