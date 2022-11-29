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
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/facebookincubator/go-belt"
	"github.com/facebookincubator/go-belt/pkg/field"
	"github.com/facebookincubator/go-belt/tool/experimental/tracer"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
)

func init() {
	if tracer.Default == nil {
		tracer.Default = func() tracer.Tracer {
			return Default()
		}
	}
}

// SpanImpl is the implementation of tracer.Span on top of zipkin.Span.
type SpanImpl struct {
	sendOnce   sync.Once
	parent     tracer.Span
	name       string
	tracer     *TracerImpl
	zipkinSpan zipkin.Span
	fields     *field.FieldsChain
}

var _ tracer.Span = (*SpanImpl)(nil)

func (t *TracerImpl) newSpanBelt(
	belt *belt.Belt,
	name string,
	parent tracer.Span,
	options ...tracer.SpanOption,
) (tracer.Span, *belt.Belt) {
	if tracer.IsNoopSpan(parent) {
		return tracer.NewNoopSpan(name, parent, time.Now()), belt
	}

	span := &SpanImpl{
		parent: parent,
		tracer: t,
		fields: t.ContextFields,
		name:   name,
	}

	if !t.PreHooks.ProcessSpan(span) {
		return tracer.NewNoopSpan(name, parent, time.Now()), belt
	}

	t.compileFields()

	zipkinOptions := spanZipkinOptions(t.compiledFields, options...)
	span.zipkinSpan = t.ZipkinTracer.StartSpan(name, zipkinOptions...)

	var returnSpan tracer.Span = span
	if zipkin.IsNoop(span.zipkinSpan) {
		returnSpan = tracer.NewNoopSpan(name, parent, time.Now())
	}

	if belt == nil {
		return returnSpan, nil
	}
	return returnSpan, tracer.BeltWithSpan(belt, span)
}

func (t *TracerImpl) newSpanCtx(
	ctx context.Context,
	name string,
	parent tracer.Span,
	options ...tracer.SpanOption,
) (tracer.Span, context.Context) {
	if tracer.IsNoopSpan(parent) {
		return tracer.NewNoopSpan(name, parent, time.Now()), ctx
	}

	span := &SpanImpl{
		parent: parent,
		tracer: t,
		fields: t.ContextFields,
		name:   name,
	}

	if !t.PreHooks.ProcessSpan(span) {
		return tracer.NewNoopSpan(name, parent, time.Now()), ctx
	}

	t.compileFields()

	zipkinOptions := spanZipkinOptions(t.compiledFields, options...)
	var resultCtx context.Context
	if ctx != nil {
		span.zipkinSpan, resultCtx = t.ZipkinTracer.StartSpanFromContext(ctx, name, zipkinOptions...)
	} else {
		span.zipkinSpan = t.ZipkinTracer.StartSpan(name, zipkinOptions...)
		resultCtx = zipkin.NewContext(context.Background(), span.zipkinSpan)
	}

	var returnSpan tracer.Span = span
	if zipkin.IsNoop(span.zipkinSpan) {
		returnSpan = tracer.NewNoopSpan(name, parent, time.Now())
	}
	resultCtx = tracer.CtxWithSpan(resultCtx, returnSpan)

	return returnSpan, resultCtx
}

func spanZipkinOptions(fields map[string]string, opts ...tracer.SpanOption) []zipkin.SpanOption {
	lastResetFieldsIdx := -1
	for idx, opt := range opts {
		if _, isResetFields := opt.(tracer.SpanOptionResetFields); isResetFields {
			lastResetFieldsIdx = idx
			continue
		}
	}
	zipkinOptions := make([]zipkin.SpanOption, 0, len(opts)+1)
	if fields != nil && lastResetFieldsIdx < 0 {
		zipkinOptions = append(zipkinOptions, zipkin.Tags(fields))
	}
	for idx, opt := range opts {
		if _, isAddFields := opt.(tracer.SpanOptionAddFields); isAddFields && idx < lastResetFieldsIdx {
			continue
		}
		zipkinOption := spanOptionToZipkin(opt)
		if zipkinOption == nil {
			continue
		}
		zipkinOptions = append(zipkinOptions, zipkinOption)
	}

	return zipkinOptions
}

func spanOptionToZipkin(opt tracer.SpanOption) zipkin.SpanOption {
	switch opt := opt.(type) {
	case tracer.SpanOptionAddFields:
		tags := make(map[string]string, len(opt))
		field.Fields(opt).ForEachField(func(f *field.Field) bool {
			tags[f.Key] = fmt.Sprint(f.Value)
			return true
		})
		return zipkin.Tags(tags)
	case tracer.SpanOptionResetFields:
		return nil
	case tracer.SpanOptionRole:
		switch opt {
		case tracer.RoleClient:
			return zipkin.Kind(model.Client)
		case tracer.RoleServer:
			return zipkin.Kind(model.Server)
		default:
			return zipkin.Kind(model.Kind(strings.ToUpper(string(opt))))
		}
	case tracer.SpanOptionStart:
		return zipkin.StartTime(time.Time(opt))
	}
	return nil
}

// ID implements tracer.Span.
func (span *SpanImpl) ID() any {
	return span.zipkinSpan.Context().ID
}

// Name implements tracer.Span.
func (span *SpanImpl) Name() string {
	return span.name
}

// TraceIDs implements tracer.Span.
func (span *SpanImpl) TraceIDs() belt.TraceIDs {
	return span.tracer.TraceIDs
}

// StartTS implements tracer.Span.
func (span *SpanImpl) StartTS() time.Time {
	b, err := json.Marshal(span.zipkinSpan)
	if err != nil {
		return time.Time{}
	}

	var spanModel model.SpanModel
	err = json.Unmarshal(b, &spanModel)
	if err != nil {
		return time.Time{}
	}

	return spanModel.Timestamp
}

// Fields implements tracer.Span.
func (span *SpanImpl) Fields() field.AbstractFields {
	return span.fields
}

// Parent implements tracer.Span.
func (span *SpanImpl) Parent() tracer.Span {
	return span.parent
}

// SetName implements tracer.Span.
func (span *SpanImpl) SetName(name string) {
	span.zipkinSpan.SetName(name)
}

// Annotate implements tracer.Span.
func (span *SpanImpl) Annotate(ts time.Time, event string) {
	span.zipkinSpan.Annotate(ts, event)
}

// SetField implements tracer.Span.
func (span *SpanImpl) SetField(k field.Key, v field.Value) {
	span.setFieldInZipkin(k, v)
}

func (span *SpanImpl) setFieldInZipkin(k field.Key, v field.Value) {
	span.fields = span.fields.WithField(k, v)
	span.zipkinSpan.Tag(k, fmt.Sprint(v))
}

// SetFields implements tracer.Span.
func (span *SpanImpl) SetFields(fields field.AbstractFields) {
	span.fields = span.fields.WithFields(fields)
	fields.ForEachField(func(f *field.Field) bool {
		span.setFieldInZipkin(f.Key, f.Value)
		return true
	})
}

// Finish implements tracer.Span.
func (span *SpanImpl) Finish() {
	span.sendOnce.Do(func() {
		if !span.tracer.Hooks.ProcessSpan(span) {
			return
		}
		span.zipkinSpan.Finish()
	})
}

// FinishWithDuration implements tracer.Span.
func (span *SpanImpl) FinishWithDuration(duration time.Duration) {
	span.sendOnce.Do(func() {
		if !span.tracer.Hooks.ProcessSpan(span) {
			return
		}
		span.zipkinSpan.FinishedWithDuration(duration)
	})
}

// SendAsIs implements tracer.Span.
func (span *SpanImpl) SendAsIs() {
	span.sendOnce.Do(span.zipkinSpan.Flush)
}

// Flush implements tracer.Span.
func (*SpanImpl) Flush() {}
