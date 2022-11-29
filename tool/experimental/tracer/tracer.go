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

package tracer

import (
	"context"

	"github.com/facebookincubator/go-belt"
	"github.com/facebookincubator/go-belt/pkg/field"
)

// Tracer is a generic abstract distributed tracing system client.
type Tracer interface {
	belt.Tool

	// Start creates a new Span, given its name, parent and options.
	Start(name string, parent Span, options ...SpanOption) Span

	// StartWithBelt creates a new root Span, given Belt, name and options.
	//
	// The returned Belt is a derivative of the provided one, with the Span added.
	StartWithBelt(belt *belt.Belt, name string, options ...SpanOption) (Span, *belt.Belt)

	// StartChildWithBelt creates a new child Span, given Belt, name and options.
	// The parent is extracted from the Belt. If one is not set in there then it is
	// an equivalent of StartWithBelt (a nil parent is used).
	//
	// The returned Belt is a derivative of the provided one, with the Span added.
	StartChildWithBelt(belt *belt.Belt, name string, options ...SpanOption) (Span, *belt.Belt)

	// StartWithCtx creates a new root Span, given Context, name and options.
	//
	// The returned Context is a derivative of the provided one, with the Span added.
	// Some implementations also injects a span structure with a specific key to the context.
	StartWithCtx(ctx context.Context, name string, options ...SpanOption) (Span, context.Context)

	// StartChildWithCtx creates a new child Span, given Context, name and options.
	// The parent is extracted from the Belt from the Context.
	// If one is not set in there then it is an equivalent of StartWithCtx (a nil parent is used).
	//
	// The returned Context is a derivative of the provided one, with the Span added.
	// Some implementations also injects a span structure with a specific key to the context.
	StartChildWithCtx(ctx context.Context, name string, options ...SpanOption) (Span, context.Context)

	// WithPreHooks returns a Tracer which includes/appends pre-hooks from the arguments.
	//
	// PreHook is the same as "Hook", but executed on early stages of building a Span
	// (before heavy computations).
	//
	// Special case: to reset hooks use `WithPreHooks()` (without any arguments).
	WithPreHooks(...Hook) Tracer

	// WithHooks returns a Tracer which includes/appends hooks from the arguments.
	//
	// See also description of "Hook".
	//
	// Special case: to reset hooks use `WithHooks()` (without any arguments).
	WithHooks(...Hook) Tracer
}

// WithField adds a Tracer derivative with the field added.
func WithField(t Tracer, key field.Key, value field.Value) Tracer {
	return t.WithContextFields(field.NewChainFromOne(key, value), 1).(Tracer)
}
