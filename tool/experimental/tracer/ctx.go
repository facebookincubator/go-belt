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
)

// StartChildSpanFromCtx creates a child span, given a context. The parent span is
// extracted from the Belt of the context. If one is not set, then the function
// returns a new root Span (a Span without a parent).
func StartChildSpanFromCtx(ctx context.Context, name string, options ...SpanOption) (Span, context.Context) {
	return FromCtx(ctx).StartChildWithCtx(ctx, name, options...)
}

// StartSpanFromCtx creates a new (root) span, given a context.
func StartSpanFromCtx(ctx context.Context, name string, options ...SpanOption) (Span, context.Context) {
	return FromCtx(ctx).StartWithCtx(ctx, name, options...)
}

// SpanFromCtx returns the current span, given a context. Returns a NoopSpan if
// one is not set.
func SpanFromCtx(ctx context.Context) Span {
	return SpanFromBelt(belt.CtxBelt(ctx))
}

// CtxWithSpan returns a context derivative/clone with the specified Span.
func CtxWithSpan(ctx context.Context, span Span) context.Context {
	return belt.CtxWithBelt(ctx, BeltWithSpan(belt.CtxBelt(ctx), span))
}

// FromCtx returns a Tracer, given a context.
func FromCtx(ctx context.Context) Tracer {
	return FromBelt(belt.CtxBelt(ctx))
}

// CtxWithTracer returns a context derivative/clone with the specified Tracer.
func CtxWithTracer(ctx context.Context, tracer Tracer) context.Context {
	return belt.WithTool(ctx, ToolID, tracer)
}
