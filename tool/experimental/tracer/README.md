# Disclaimer

This API is experimental and has no stability guarantees.

# Example
```go
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
```

# Interface

```go
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

	// Flush forces to flush all buffers.
	Flush()
}
```

