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

package tests

import (
	"github.com/facebookincubator/go-belt/tool/experimental/tracer"
	loggerimpl "github.com/facebookincubator/go-belt/tool/experimental/tracer/implementation/logger"
	"github.com/facebookincubator/go-belt/tool/logger"
	"github.com/facebookincubator/go-belt/tool/logger/adapter"
	"github.com/facebookincubator/go-belt/tool/logger/types"
)

type dummyLoggerReporter struct {
	onSend func(span tracer.Span)
	tracer *loggerimpl.TracerImpl
}

func (l *dummyLoggerReporter) Flush() {}
func (l *dummyLoggerReporter) Emit(entry *logger.Entry) {
	restoredSpan := loggerimpl.SpanFromLogEntry(entry)
	restoredSpan.Tracer = l.tracer
	l.onSend(restoredSpan)
}

func (l *dummyLoggerReporter) OnSend(onSend func(span tracer.Span)) {
	l.onSend = onSend
}

func init() {
	implementations = append(implementations, implementationCase{
		Name: "logger",
		Factory: func() (tracer.Tracer, DummyReporter) {
			var reporter dummyLoggerReporter
			tracer := loggerimpl.New(adapter.LoggerFromEmitter(&reporter).WithLevel(types.LevelTrace), func(span *loggerimpl.SpanImpl) types.Level {
				return types.LevelTrace
			})
			reporter.tracer = tracer
			if tracer == nil {
				panic("tracer == nil")
			}
			return tracer, &reporter
		},
	})
}
