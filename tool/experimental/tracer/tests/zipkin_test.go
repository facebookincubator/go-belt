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
	"fmt"
	"time"

	"github.com/facebookincubator/go-belt/tool/experimental/tracer"
	zipkinimpl "github.com/facebookincubator/go-belt/tool/experimental/tracer/implementation/zipkin"
	"github.com/openzipkin/zipkin-go"
	zipkinmodel "github.com/openzipkin/zipkin-go/model"
	zipkinreporter "github.com/openzipkin/zipkin-go/reporter"
	"github.com/xaionaro-go/unsafetools"
)

type dummyZipkinReporter struct {
	onSend func(tracer.Span)
	tracer *zipkinimpl.TracerImpl
}

type dummyZipkinSpan struct {
	data     zipkinmodel.SpanModel
	reporter *dummyZipkinReporter
}

func (span *dummyZipkinSpan) Context() zipkinmodel.SpanContext {
	return span.data.SpanContext
}
func (span *dummyZipkinSpan) SetName(name string) {
	span.data.Name = name
}
func (span *dummyZipkinSpan) SetRemoteEndpoint(*zipkinmodel.Endpoint) {
	panic("not implemented")
}
func (span *dummyZipkinSpan) Annotate(time.Time, string) {
	panic("not implemented")
}
func (span *dummyZipkinSpan) Tag(key, value string) {
	span.data.Tags[key] = value
}
func (span *dummyZipkinSpan) Finish() {
	panic("not implemented")
}
func (span *dummyZipkinSpan) FinishedWithDuration(duration time.Duration) {
	panic("not implemented")
}
func (span *dummyZipkinSpan) Flush() {
	panic("not implemented")
}

func (r *dummyZipkinReporter) Send(span zipkinmodel.SpanModel) {
	spanImpl := &zipkinimpl.SpanImpl{}
	*unsafetools.FieldByName(spanImpl, "name").(*string) = span.Name
	*unsafetools.FieldByName(spanImpl, "zipkinSpan").(*zipkin.Span) = &dummyZipkinSpan{data: span, reporter: r}
	*unsafetools.FieldByName(spanImpl, "tracer").(**zipkinimpl.TracerImpl) = r.tracer
	r.onSend(spanImpl)
}

func (r *dummyZipkinReporter) OnSend(onSend func(tracer.Span)) {
	r.onSend = onSend
}

func (r *dummyZipkinReporter) Close() error {
	return nil
}

var _ zipkinreporter.Reporter = (*dummyZipkinReporter)(nil)

func init() {
	implementations = append(implementations, implementationCase{
		Name: "zipkin",
		Factory: func() (tracer.Tracer, DummyReporter) {
			var reporter dummyZipkinReporter

			zipkinTracer, err := zipkin.NewTracer(&reporter)
			if err != nil {
				panic(fmt.Errorf("unable to create a zipkin tracer: %v", err))
			}

			tracer := zipkinimpl.New(zipkinTracer)
			return tracer, &reporter
		},
	})
}
