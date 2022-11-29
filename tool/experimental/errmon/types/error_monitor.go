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
)

// Emitter is a sender of already generated Events.
type Emitter interface {
	// Flush forces to flush all buffers.
	Flush()

	// SendEvent sends an Event and optionally adds a value to ExternalIDs.
	Emit(*Event)
}

// Emitters is a collection of Emitter-s.
type Emitters []Emitter

var _ Emitter = (Emitters)(nil)

// Flush implements Emitter
func (s Emitters) Flush() {
	for _, emitter := range s {
		emitter.Flush()
	}
}

// Emit implements Emitter
func (s Emitters) Emit(ev *Event) {
	for _, emitter := range s {
		emitter.Emit(ev)
	}
}

// ErrorMonitor is an observability Tool (belt.Tool) which allows
// to report about any exceptions which happen for debugging. It
// collects any useful information it can.
//
// An ErrorMonitor implementation is not supposed to be fast, but
// it supposed to provide verbose reports (sufficient enough to
// debug found problems).
type ErrorMonitor interface {
	belt.Tool

	// Emitter returns the Emitter.
	//
	// A read-only value, do not change it.
	Emitter() Emitter

	// ObserveError issues an error event if `err` is not an untyped nil. Additional
	// data (left by various observability tooling) is extracted from `belt`.
	//
	// Returns an Event only if one was issued (and for example was not sampled out by a Sampler Hook).
	ObserveError(*belt.Belt, error) *Event

	// ObserveRecover issues a panic event if `recoverResult` is not an untyped nil.
	// Additional data (left by various observability tooling) is extracted from `belt`.
	//
	// Is supposed to be used in constructions like:
	//
	//     defer func() { errmon.ObserveRecover(ctx, recover()) }()
	//
	// See also: https://go.dev/ref/spec#Handling_panics
	//
	// Returns an Event only if one was issued (and for example was not sampled out by a Sampler Hook).
	ObserveRecover(_ *belt.Belt, recoverResult any) *Event

	// WithPreHooks returns a ErrorMonitor derivative which also includes/appends pre-hooks from the arguments.
	//
	// Special case: to reset hooks use `WithHooks()` (without any arguments).
	WithPreHooks(...PreHook) ErrorMonitor

	// WithHooks returns a ErrorMonitor derivative which also includes/appends hooks from the arguments.
	//
	// Special case: to reset hooks use `WithHooks()` (without any arguments).
	WithHooks(...Hook) ErrorMonitor
}
