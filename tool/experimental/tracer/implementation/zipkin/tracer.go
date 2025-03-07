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

package zipkin

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/facebookincubator/go-belt"
	"github.com/facebookincubator/go-belt/pkg/field"
	"github.com/facebookincubator/go-belt/tool/experimental/tracer"
	"github.com/openzipkin/zipkin-go"
	logreporter "github.com/openzipkin/zipkin-go/reporter/log"
)

// TracerImpl is the implementation of tracer.Tracer on top of zipkin.Tracer.
type TracerImpl struct {
	ZipkinTracer      *zipkin.Tracer
	ContextFields     *field.FieldsChain
	compileFieldsOnce sync.Once
	compiledFields    map[string]string
	TraceIDs          belt.TraceIDs
	TagTraceIDs       string
	Hooks             tracer.Hooks
	PreHooks          tracer.Hooks
}

var _ tracer.Tracer = (*TracerImpl)(nil)

func (t TracerImpl) clone() *TracerImpl { //nolint:govet
	t.compileFieldsOnce = sync.Once{}
	return &t
}

// WithContextFields implements tracer.Tracer.
func (t *TracerImpl) WithContextFields(allFields *field.FieldsChain, newFieldsCount int) belt.Tool {
	c := t.clone()
	c.ContextFields = allFields
	return c
}

// WithTraceIDs implements tracer.Tracer.
func (t *TracerImpl) WithTraceIDs(traceIDs belt.TraceIDs, newTraceIDsCount int) belt.Tool {
	c := t.clone()
	c.TraceIDs = traceIDs
	return c
}

// Start implements tracer.Tracer.
func (t *TracerImpl) Start(name string, parent tracer.Span, options ...tracer.SpanOption) tracer.Span {
	span, _ := t.newSpanBelt(nil, name, parent, options...)
	return span
}

// StartWithCtx implements tracer.Tracer.
func (t *TracerImpl) StartWithCtx(ctx context.Context, name string, options ...tracer.SpanOption) (tracer.Span, context.Context) {
	return t.newSpanCtx(ctx, name, nil, options...)

}

// StartChildWithCtx implements tracer.Tracer.
func (t *TracerImpl) StartChildWithCtx(ctx context.Context, name string, options ...tracer.SpanOption) (tracer.Span, context.Context) {
	return t.newSpanCtx(ctx, name, tracer.SpanFromCtx(ctx), options...)
}

// StartWithBelt implements tracer.Tracer.
func (t *TracerImpl) StartWithBelt(belt *belt.Belt, name string, options ...tracer.SpanOption) (tracer.Span, *belt.Belt) {
	return t.newSpanBelt(belt, name, nil, options...)

}

// StartChildWithBelt implements tracer.Tracer.
func (t *TracerImpl) StartChildWithBelt(belt *belt.Belt, name string, options ...tracer.SpanOption) (tracer.Span, *belt.Belt) {
	return t.newSpanBelt(belt, name, tracer.SpanFromBelt(belt), options...)
}

// DefaultTagTraceIDs is the tag name used to store belt.TraceIDs value.
var DefaultTagTraceIDs = `trace_ids`

func (t *TracerImpl) compileFields() {
	t.compileFieldsOnce.Do(func() {
		t.compiledFields = make(map[string]string, t.ContextFields.Len()+1)
		t.ContextFields.ForEachField(func(f *field.Field) bool {
			t.compiledFields[f.Key] = fmt.Sprint(f.Value)
			return true
		})
		if t.TraceIDs == nil {
			return
		}
		strs := make([]string, 0, len(t.TraceIDs))
		for _, traceID := range t.TraceIDs {
			strs = append(strs, string(traceID))
		}
		t.compiledFields[t.TagTraceIDs] = strings.Join(strs, ",")
	})
}

// WithPreHooks implements tracer.Tracer.
func (t *TracerImpl) WithPreHooks(hooks ...tracer.Hook) tracer.Tracer {
	c := t.clone()
	if hooks == nil {
		c.PreHooks = nil
	} else {
		c.PreHooks = tracer.Hooks{c.PreHooks, tracer.Hooks(hooks)}
	}
	return c
}

// WithHooks implements tracer.Tracer.
func (t *TracerImpl) WithHooks(hooks ...tracer.Hook) tracer.Tracer {
	c := t.clone()
	if hooks == nil {
		c.Hooks = nil
	} else {
		c.Hooks = tracer.Hooks{c.Hooks, tracer.Hooks(hooks)}
	}
	return c
}

// Flush implements tracer.Tracer.
func (t *TracerImpl) Flush(context.Context) {
	panic("not supported")
}

// Default returns the default tracer.Tracer on top of zipkin.
var Default = func() tracer.Tracer {
	reporter := logreporter.NewReporter(log.Default())
	defer reporter.Close()

	tracer, err := zipkin.NewTracer(reporter)
	if err != nil {
		log.Fatalf("unable to create a zipkin tracer: %v", err)
	}

	return New(tracer)
}

// New returns a new instance of TracerImpl.
func New(zipkinTracer *zipkin.Tracer, options ...Option) *TracerImpl {
	tracer := &TracerImpl{
		ZipkinTracer: zipkinTracer,
		TagTraceIDs:  DefaultTagTraceIDs,
	}
	for _, opt := range options {
		opt.apply(tracer)
	}
	return tracer
}
