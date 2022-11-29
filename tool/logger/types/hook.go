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

package types

import (
	"github.com/facebookincubator/go-belt"
	"github.com/facebookincubator/go-belt/pkg/field"
)

// Hook is a pre-processor for log entries.
// It may for example modify them or cancel them.
type Hook interface {
	// ProcessLogEntry is called right before sending the Entry to a Emitter.
	//
	// If the returned value is false, then the sending and processing by other hooks
	// are cancelled.
	//
	// It is allowed to modify the entry.
	ProcessLogEntry(*Entry) bool

	// Flush gracefully empties any buffers the hook may have.
	Flush()
}

// Hooks is a collection of Hook-s.
type Hooks []Hook

var _ Hook = (Hooks)(nil)

// ProcessLogEntry implements Hook.
func (s Hooks) ProcessLogEntry(e *Entry) bool {
	for _, hook := range s {
		if !hook.ProcessLogEntry(e) {
			return false
		}
	}
	return true
}

// Flush implements Hook.
func (s Hooks) Flush() {
	for _, hook := range s {
		hook.Flush()
	}
}

// PreHookResult is a result of computations returned by a PreHook.
type PreHookResult struct {
	// Skip forces a Logger to do not log this Entry.
	Skip bool

	// ExtraFields adds additional fields to the Entry.
	ExtraFields field.AbstractFields

	// ExtraEntryProperties adds additional Entry.Properties.
	ExtraEntryProperties EntryProperties
}

// PreHook is executed very early in a Logger, allowing to discard or change
// an Entry before any essential computations are done.
//
// For example it may be useful for samplers.
type PreHook interface {
	// ProcessInput is executed when functions Log, Trace, Debug, Info, Warn, Error, Panic and Fatal are called.
	//
	// TraceIDs are provided by Logger/CompactLogger and the rest arguments are just passed-through.
	ProcessInput(belt.TraceIDs, Level, ...any) PreHookResult

	// ProcessInputf is executed when functions Logf, Tracef, Debugf, Infof, Warnf, Errorf, Panicf and Fatalf are called.
	//
	// TraceIDs are provided by Logger/CompactLogger and the rest arguments are just passed-through.
	ProcessInputf(belt.TraceIDs, Level, string, ...any) PreHookResult

	// ProcessInputf is executed when functions LogFields, TraceFields, DebugFields, InfoFields, WarnFields, ErrorFields, PanicFields and FatalFields are called.
	//
	// TraceIDs are provided by Logger/CompactLogger and the rest arguments are just passed-through.
	ProcessInputFields(belt.TraceIDs, Level, string, field.AbstractFields) PreHookResult
}

// PreHooks is a collection of PreHook-s.
type PreHooks []PreHook

var _ PreHook = (PreHooks)(nil)

// ProcessInput implements PreHook.
func (s PreHooks) ProcessInput(traceIDs belt.TraceIDs, level Level, values ...any) PreHookResult {
	var result PreHookResult
	for _, hook := range s {
		oneResult := hook.ProcessInput(traceIDs, level, values...)
		result.Skip = oneResult.Skip
		if result.Skip {
			return result
		}
		result.ExtraEntryProperties = append(result.ExtraEntryProperties, oneResult.ExtraEntryProperties)
		result.ExtraFields = field.Add(result.ExtraFields, oneResult.ExtraFields)
	}

	return result
}

// ProcessInputf implements PreHook.
func (s PreHooks) ProcessInputf(traceIDs belt.TraceIDs, level Level, format string, args ...any) PreHookResult {
	var result PreHookResult
	for _, hook := range s {
		oneResult := hook.ProcessInputf(traceIDs, level, format, args...)
		result.Skip = oneResult.Skip
		if result.Skip {
			return result
		}
		result.ExtraEntryProperties = append(result.ExtraEntryProperties, oneResult.ExtraEntryProperties)
		result.ExtraFields = field.Add(result.ExtraFields, oneResult.ExtraFields)
	}

	return result
}

// ProcessInputFields implements PreHook.
func (s PreHooks) ProcessInputFields(traceIDs belt.TraceIDs, level Level, message string, addFields field.AbstractFields) PreHookResult {
	var result PreHookResult
	for _, hook := range s {
		oneResult := hook.ProcessInputFields(traceIDs, level, message, addFields)
		result.Skip = oneResult.Skip
		if result.Skip {
			return result
		}
		result.ExtraEntryProperties = append(result.ExtraEntryProperties, oneResult.ExtraEntryProperties)
		result.ExtraFields = field.Add(result.ExtraFields, oneResult.ExtraFields)
	}

	return result
}
