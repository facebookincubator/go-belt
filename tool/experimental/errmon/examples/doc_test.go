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
	errmon.Default = func(b *belt.Belt) errmon.ErrorMonitor {
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
	m := errmon.Default(nil)
	m = errmon.WithField(m, "arg", arg)
	defer func() { m.ObserveRecover(nil, recover()) }()
	// ..do something here..
}
