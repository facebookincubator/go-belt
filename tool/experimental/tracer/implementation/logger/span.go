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

package logger

import (
	"context"
	"time"

	"github.com/facebookincubator/go-belt"
	"github.com/facebookincubator/go-belt/pkg/field"
	"github.com/facebookincubator/go-belt/tool/experimental/tracer"
)

// Event is an unit of information added through method tracer.Span.Annotate.
type Event struct {
	Timestamp time.Time
	Name      string
}

// SpanImpl is the implementation of tracer.Span.
type SpanImpl struct {
	Tracer       *TracerImpl
	StartTSValue time.Time
	Duration     time.Duration
	ParentValue  tracer.Span
	FieldsValue  *field.FieldsChain
	NameValue    string
	Events       []Event
}

var _ tracer.Span = (*SpanImpl)(nil)

func (t *TracerImpl) newSpanCtx(
	ctx context.Context,
	name string,
	parent tracer.Span,
	options ...tracer.SpanOption,
) (tracer.Span, context.Context) {
	span := &SpanImpl{
		Tracer:       t,
		StartTSValue: time.Now(),
		ParentValue:  parent,
		FieldsValue:  t.Fields,
		NameValue:    name,
	}
	if !t.PreHooks.ProcessSpan(span) {
		return tracer.NewNoopSpan(name, parent, time.Now()), ctx
	}

	var ctxDer context.Context
	if ctx != nil {
		ctxDer = tracer.CtxWithSpan(ctx, span)
	}

	return span, ctxDer
}

func (t *TracerImpl) newSpanBelt(
	_belt *belt.Belt,
	name string,
	parent tracer.Span,
	options ...tracer.SpanOption,
) (tracer.Span, *belt.Belt) {
	span := &SpanImpl{
		Tracer:       t,
		StartTSValue: time.Now(),
		ParentValue:  parent,
		FieldsValue:  t.Fields,
		NameValue:    name,
	}
	if !t.PreHooks.ProcessSpan(span) {
		return tracer.NewNoopSpan(name, parent, time.Now()), _belt
	}

	var beltDer *belt.Belt
	if _belt != nil {
		beltDer = tracer.BeltWithSpan(_belt, span)
	}

	return span, beltDer
}

// ID implements tracer.Span
func (span *SpanImpl) ID() any {
	return span.TraceIDs()
}

// TraceIDs implements tracer.Span
func (span *SpanImpl) TraceIDs() belt.TraceIDs {
	return span.Tracer.TraceIDs
}

// Name implements tracer.Span
func (span *SpanImpl) Name() string {
	return span.NameValue
}

// StartTS implements tracer.Span
func (span *SpanImpl) StartTS() time.Time {
	return span.StartTSValue
}

// Fields implements tracer.Span
func (span *SpanImpl) Fields() field.AbstractFields {
	return span.FieldsValue
}

// Parent implements tracer.Span
func (span *SpanImpl) Parent() tracer.Span {
	return span.ParentValue
}

// SetName implements tracer.Span
func (span *SpanImpl) SetName(name string) {
	span.NameValue = name
}

// Annotate implements tracer.Span
func (span *SpanImpl) Annotate(ts time.Time, event string) {
	span.Events = append(span.Events, Event{
		Timestamp: ts,
		Name:      event,
	})
}

// SetField implements tracer.Span
func (span *SpanImpl) SetField(key field.Key, value field.Value) {
	span.FieldsValue = span.FieldsValue.WithField(key, value)
}

// SetFields implements tracer.Span
func (span *SpanImpl) SetFields(fields field.AbstractFields) {
	span.FieldsValue = span.FieldsValue.WithFields(fields)
}

// Finish implements tracer.Span
func (span *SpanImpl) Finish() {
	span.Duration = time.Since(span.StartTSValue)
	span.Tracer.send(span)
}

// FinishWithDuration implements tracer.Span
func (span *SpanImpl) FinishWithDuration(duration time.Duration) {
	span.Duration = duration
	span.Tracer.send(span)
}

// Flush implements tracer.Span
func (span *SpanImpl) Flush() {
	span.Tracer.Flush()
}
