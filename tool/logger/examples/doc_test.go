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

package examples

import (
	"bytes"
	"context"
	"log"

	"github.com/facebookincubator/go-belt"
	"github.com/facebookincubator/go-belt/tool/logger"
	"github.com/facebookincubator/go-belt/tool/logger/implementation/glog"
	xlogrus "github.com/facebookincubator/go-belt/tool/logger/implementation/logrus"
	"github.com/facebookincubator/go-belt/tool/logger/implementation/stdlib"
	"github.com/facebookincubator/go-belt/tool/logger/implementation/zap"
	"github.com/sirupsen/logrus"
)

func Example() {
	// easy to use:
	l := xlogrus.Default()
	someFunction(1, l)

	// implementation agnostic:
	l = zap.Default()
	someFunction(2, l)

	// various implementations:
	l = glog.New()
	someFunction(3, l)

	// one may still reuse all the features of the backend logger:
	logrusInstance := logrus.New()
	logrusInstance.Formatter = &logrus.JSONFormatter{}
	l = xlogrus.New(logrusInstance)
	someFunction(4, l)

	// just another example:
	var buf bytes.Buffer
	stdLogInstance := log.New(&buf, "", log.Llongfile)
	l = stdlib.New(stdLogInstance, logger.LevelDebug)
	someFunction(5, l)

	// use go-belt to manage the logger
	obs := belt.New()
	obs = logger.BeltWithLogger(obs, l)
	someBeltyFunction(6, obs)

	// use context to manage the logger
	ctx := context.Background()
	ctx = logger.CtxWithLogger(ctx, l)
	someContextyFunction(ctx, 7)

	// use a singletony Logger:
	logger.Default = func() logger.Logger {
		return l
	}
	oneMoreFunction(8)
}

func someFunction(arg int, l logger.Logger) {
	// experience close to logrus/zap:
	l = l.WithField("arg", arg)
	anotherFunction(l)
}

func anotherFunction(l logger.Logger) {
	l.Debugf("hello world, %T!", l)
}

func someBeltyFunction(arg int, obs *belt.Belt) {
	obs = obs.WithField("arg", arg)
	anotherBeltyFunction(obs)
}

func anotherBeltyFunction(obs *belt.Belt) {
	logger.FromBelt(obs).Debugf("hello world!")
}

func someContextyFunction(ctx context.Context, arg int) {
	ctx = belt.WithField(ctx, "arg", arg)
	anotherContextyFunction(ctx)
}

func anotherContextyFunction(ctx context.Context) {
	logger.FromCtx(ctx).Debugf("hello world!")
	// or a shorter form:
	logger.Debugf(ctx, "hello world!")
}

func oneMoreFunction(arg int) {
	logger.Default().WithField("arg", arg).Debugf("hello world!")
}
