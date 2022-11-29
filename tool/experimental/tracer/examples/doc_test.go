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

package tracer_test

import (
	"context"

	"github.com/facebookincubator/go-belt"
	"github.com/facebookincubator/go-belt/beltctx"
	"github.com/facebookincubator/go-belt/tool/experimental/tracer"
	"github.com/facebookincubator/go-belt/tool/experimental/tracer/implementation/logger"
	"github.com/facebookincubator/go-belt/tool/experimental/tracer/implementation/zipkin"
)

func Example() {
	// easy to use:
	t := logger.Default()
	someFunction(1, t)

	// implementation agnostic:
	t = zipkin.Default()
	someFunction(2, t)

	// one may still reuse all the features of the backend Tracer:
	t.(*zipkin.TracerImpl).ZipkinTracer.SetNoop(true)

	// use go-belt to manage the Tracer
	obs := belt.New()
	obs = tracer.BeltWithTracer(obs, t)
	someBeltyFunction(3, obs)

	// use context to manage the Tracer
	ctx := context.Background()
	ctx = tracer.CtxWithTracer(ctx, t)
	someContextyFunction(ctx, 4)

	// use a singletony Tracer:
	tracer.Default = func() tracer.Tracer {
		return t
	}
	yetOneMoreFunction(5)
}

func someFunction(arg int, t tracer.Tracer) {
	// experience close to logrus/zap:
	t = tracer.WithField(t, "arg", arg)
	anotherFunction(t)
}

func anotherFunction(t tracer.Tracer) {
	span := t.Start("hello", nil)
	defer span.Finish()
	// ..do something long here..
	oneMoreFunction(t, span)
}

func oneMoreFunction(t tracer.Tracer, parentSpan tracer.Span) {
	span := t.Start("child", parentSpan)
	defer span.Finish()
	// ..do something meaningful here..
}

func someBeltyFunction(arg int, obs *belt.Belt) {
	obs = obs.WithField("arg", arg)
	anotherBeltyFunction(obs)
}

func anotherBeltyFunction(obs *belt.Belt) {
	span, obs := tracer.StartChildSpanFromBelt(obs, "hello")
	defer span.Finish()
	// ..do something long here..
	oneMoreBeltyFunction(obs)
}

func oneMoreBeltyFunction(obs *belt.Belt) {
	span, obs := tracer.StartChildSpanFromBelt(obs, "child")
	defer span.Finish()
	// ..do something meaningful here..
	_ = obs
}

func someContextyFunction(ctx context.Context, arg int) {
	ctx = beltctx.WithField(ctx, "arg", arg)
	anotherContextyFunction(ctx)
}

func anotherContextyFunction(ctx context.Context) {
	span, ctx := tracer.StartChildSpanFromCtx(ctx, "hello")
	defer span.Finish()
	// ..do something long here..
	oneMoreContextyFunction(ctx)
}

func oneMoreContextyFunction(ctx context.Context) {
	span, ctx := tracer.StartChildSpanFromCtx(ctx, "child")
	defer span.Finish()
	// ..do something meaningful here..
	_ = ctx
}

func yetOneMoreFunction(arg int) {
	t := tracer.Default()
	t = tracer.WithField(t, "arg", arg)
	span := t.Start("hello", nil)
	defer span.Finish()
}
