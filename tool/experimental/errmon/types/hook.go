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

// PreHookResult is the result of a PreHook.
type PreHookResult struct {
	Skip        bool
	ExtraFields field.AbstractFields
}

// PreHook is similar to a Hook, but is used before all the information is collected
// and the Event is generated. It might be useful for example for a Sampler to avoid
// computations for sampled out events.
type PreHook interface {
	ProcessInputError(belt.TraceIDs, error) PreHookResult
	ProcessInputPanic(belt.TraceIDs, any) PreHookResult
}

// PreHooks is a collection of PreHook-s.
type PreHooks []PreHook

var _ PreHook = (PreHooks)(nil)

// ProcessInputError implements PreHook.
func (s PreHooks) ProcessInputError(traceIDs belt.TraceIDs, err error) PreHookResult {
	var result PreHookResult
	for _, hook := range s {
		oneResult := hook.ProcessInputError(traceIDs, err)
		result.Skip = oneResult.Skip
		if result.Skip {
			return result
		}
		result.ExtraFields = field.Add(result.ExtraFields, oneResult.ExtraFields)
	}

	return result
}

// ProcessInputPanic implements PreHook.
func (s PreHooks) ProcessInputPanic(traceIDs belt.TraceIDs, panicValue any) PreHookResult {
	var result PreHookResult
	for _, hook := range s {
		oneResult := hook.ProcessInputPanic(traceIDs, panicValue)
		result.Skip = oneResult.Skip
		if result.Skip {
			return result
		}
		result.ExtraFields = field.Add(result.ExtraFields, oneResult.ExtraFields)
	}

	return result
}

// Hook is a pre-processor for an Event which is ran before sending it.
// It may modify the event or prevent from being sent.
type Hook interface {
	// Process performs the modification of an Event and/or returns
	// false to prevent an Event from being sent (or returns true
	// to allow sending the Event).
	Process(*Event) bool
}

// Hooks is a collection of Hook-s
type Hooks []Hook

var _ Hook = (Hooks)(nil)

// Process implements Hook.
func (s Hooks) Process(ev *Event) bool {
	for _, hook := range s {
		if !hook.Process(ev) {
			return false
		}
	}
	return true
}
