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

// CompactLogger is an optimized interface of an abstract
// generic structured logger. For the full interface see "Logger".
//
// It is easier to implement a CompactLogger and then wrap it
// to an adapter.GenericSugar to get a Logger, than to implement
// a full Logger.
//
// All methods are thread-safe.
type CompactLogger interface {
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
	Logf(level types.Level, format string, args ...any)

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
	LogFields(level types.Level, message string, fields field.AbstractFields)

	// Log extracts structured fields from provided values, joins
	// the rest into an unstructured message and logs the result.
	//
	// This function provides convenience (relatively to LogFields)
	// at cost of a bit of performance.
	//
	// There are few ways to extract structured fields, which are
	// applied for each value from `values` (in descending priority order):
	// 1. If a `value` is an `*Entry` then the Entry is used (with its fields)
	// 2. If a `value` implements field.AbstractFields then ForEachField method
	//    is used (so it is become similar to LogFields).
	// 3. If a `value` is a structure (or a pointer to a structure) then
	//    fields of the structure are interpreted as structured fields
	//    to be logged (see explanation below).
	// 4. If a `value` is a map then fields a constructed out of this map.
	//
	// Everything that does not fit into any of the rules above is just
	// joined into an unstructured message (and works the same way
	// as `message` in LogFields).
	//
	// How structures are parsed:
	// Structures are parsed recursively. Each field name of the path in a tree
	// of structures is added to the resulting field name (for example int "struct{A struct{B int}}"
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
	Log(level types.Level, values ...any)

	// Emitter returns the Emitter (see the description of interface "Emitter").
	Emitter() types.Emitter

	// Level returns the current logging level (see description of "Level").
	Level() types.Level

	// WithLevel returns a logger with logger level set to the given argument.
	//
	// See also the description of type "Level".
	WithLevel(types.Level) CompactLogger

	// WithPreHooks returns a Logger which includes/appends pre-hooks from the arguments.
	//
	// See also description of "PreHook".
	//
	// Special case: to reset hooks use `WithHooks()` (without any arguments).
	WithPreHooks(...types.PreHook) CompactLogger

	// WithHooks returns a Logger which includes/appends hooks from the arguments.
	//
	// See also description of "Hook".
	//
	// Special case: to reset hooks use `WithHooks()` (without any arguments).
	WithHooks(...types.Hook) CompactLogger

	// WithField returns the logger with the added field (used for structured logging).
	WithField(key string, value any, props ...field.Property) CompactLogger

	// WithFields returns the logger with the fields added (used for structured logging).
	WithFields(fields field.AbstractFields) CompactLogger

	// WithMessagePrefix adds a string to all messages logged through the derived logger
	WithMessagePrefix(prefix string) CompactLogger

	// WithEntryProperties adds props to EntryProperties of each emitted Entry.
	// This could be used only for enabling implementation-specific behavior.
	WithEntryProperties(props ...types.EntryProperty) CompactLogger
}
