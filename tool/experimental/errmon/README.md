# Disclaimer

This API is experimental and has no stability guarantees.

# Example
```go
package errmon_test

import (
	"context"

	sentryupstream "github.com/getsentry/sentry-go"

	"github.com/facebookincubator/go-belt"
	"github.com/facebookincubator/go-belt/beltctx"
	"github.com/facebookincubator/go-belt/tool/experimental/errmon"
	"github.com/facebookincubator/go-belt/tool/experimental/errmon/implementation/logger"
	"github.com/facebookincubator/go-belt/tool/experimental/errmon/implementation/sentry"
)

func Example() {
	// easy to use:
	m := logger.Default()
	someFunction(1, m)

	// implementation agnostic:
	sentryClient, err := sentryupstream.NewClient(sentryupstream.ClientOptions{})
	if err != nil {
		panic(err)
	}
	m = sentry.New(sentryClient)
	someFunction(2, m)

	// one may still reuse all the features of the emitter ErrorMonitor:
	_ = m.Emitter().(*sentry.Emitter).SentryClient.Options()

	// use go-belt to manage the ErrorMonitor
	obs := belt.New()
	obs = errmon.BeltWithErrorMonitor(obs, m)
	someBeltyFunction(3, obs)

	// use context to manage the ErrorMonitor
	ctx := context.Background()
	ctx = errmon.CtxWithErrorMonitor(ctx, m)
	someContextyFunction(ctx, 4)

	// use a singletony ErrorMonitor:
	errmon.Default = func() errmon.ErrorMonitor {
		return m
	}
	yetOneMoreFunction(5)
}

func someFunction(arg int, m errmon.ErrorMonitor) {
	// experience close to logrus/zap:
	m = errmon.WithField(m, "arg", arg)
	anotherFunction(m)
}

func anotherFunction(m errmon.ErrorMonitor) {
	defer func() { m.ObserveRecover(nil, recover()) }()
	// ..do something here..
}

func someBeltyFunction(arg int, obs *belt.Belt) {
	obs = obs.WithField("arg", arg)
	anotherBeltyFunction(obs)
}

func anotherBeltyFunction(obs *belt.Belt) {
	defer func() { errmon.ObserveRecoverBelt(obs, recover()) }()
	// ..do something here..
}

func someContextyFunction(ctx context.Context, arg int) {
	ctx = beltctx.WithField(ctx, "arg", arg)
	anotherContextyFunction(ctx)
}

func anotherContextyFunction(ctx context.Context) {
	defer func() { errmon.ObserveRecoverCtx(ctx, recover()) }()
	// ..do something here..
}

func yetOneMoreFunction(arg int) {
	m := errmon.Default()
	m = errmon.WithField(m, "arg", arg)
	defer func() { m.ObserveRecover(nil, recover()) }()
	// ..do something here..
}
```

# Interface
```go
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
```

