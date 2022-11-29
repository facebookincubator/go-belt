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
	"time"

	"github.com/facebookincubator/go-belt"
	"github.com/facebookincubator/go-belt/pkg/field"
)

// NoopSpan is a no-op implementation of a Tracer. Supposed
// to be used as a dummy placeholder if no Tracer is set.
type noopTracer struct{}

// WithContextFields implements Tracer.
func (t noopTracer) WithContextFields(allFields *field.FieldsChain, newFieldsCount int) belt.Tool {
	return t
}

// WithTraceIDs implements Tracer.
func (t noopTracer) WithTraceIDs(traceIDs belt.TraceIDs, newTraceIDsCount int) belt.Tool {
	return t
}

// Start implements Tracer.
func (noopTracer) Start(name string, parent Span, options ...SpanOption) Span {
	return NewNoopSpan(name, parent, time.Time{})
}

// StartWithBelt implements Tracer.
func (noopTracer) StartWithBelt(belt *belt.Belt, name string, options ...SpanOption) (Span, *belt.Belt) {
	span := NewNoopSpan(name, nil, time.Time{})
	return span, BeltWithSpan(belt, span)
}

// StartChildWithBelt implements Tracer.
func (noopTracer) StartChildWithBelt(belt *belt.Belt, name string, options ...SpanOption) (Span, *belt.Belt) {
	span := NewNoopSpan(name, SpanFromBelt(belt), time.Time{})
	return span, BeltWithSpan(belt, span)
}

// StartWithCtx implements Tracer.
func (noopTracer) StartWithCtx(ctx context.Context, name string, options ...SpanOption) (Span, context.Context) {
	span := NewNoopSpan(name, nil, time.Time{})
	return span, CtxWithSpan(ctx, span)
}

// StartChildWithCtx implements Tracer.
func (noopTracer) StartChildWithCtx(ctx context.Context, name string, options ...SpanOption) (Span, context.Context) {
	span := NewNoopSpan(name, SpanFromCtx(ctx), time.Time{})
	return span, CtxWithSpan(ctx, span)
}

// WithPreHooks implements Tracer.
func (t noopTracer) WithPreHooks(...Hook) Tracer {
	return t
}

// WithHooks implements Tracer.
func (t noopTracer) WithHooks(...Hook) Tracer {
	return t
}

// Flush implements Tracer.
func (noopTracer) Flush() {}
