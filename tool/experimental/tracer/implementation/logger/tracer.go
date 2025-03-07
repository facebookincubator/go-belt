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
	"github.com/facebookincubator/go-belt/tool/logger"
	loggertypes "github.com/facebookincubator/go-belt/tool/logger/types"
)

// TracerImpl is the implementation of tracer.Tracer based on a given logger.Logger.
type TracerImpl struct {
	Logger    loggertypes.Logger
	LevelFunc func(span *SpanImpl) loggertypes.Level
	PreHooks  tracer.Hooks
	Hooks     tracer.Hooks
	TraceIDs  belt.TraceIDs
	Fields    *field.FieldsChain
}

var _ tracer.Tracer = (*TracerImpl)(nil)

func (t TracerImpl) clone() *TracerImpl {
	return &t
}

// WithContextFields implements tracer.Tracer.
func (t *TracerImpl) WithContextFields(allFields *field.FieldsChain, newFieldsCount int) belt.Tool {
	clone := t.clone()
	clone.Fields = allFields
	return clone
}

// WithTraceIDs implements tracer.Tracer.
func (t *TracerImpl) WithTraceIDs(traceIDs belt.TraceIDs, newTraceIDsCount int) belt.Tool {
	clone := t.clone()
	clone.TraceIDs = traceIDs
	clone.Logger = t.Logger.WithTraceIDs(traceIDs, newTraceIDsCount).(logger.Logger)
	return clone
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

// Start implements tracer.Tracer.
func (t *TracerImpl) Start(name string, parent tracer.Span, options ...tracer.SpanOption) tracer.Span {
	span, _ := t.newSpanBelt(nil, name, parent, options...)
	return span
}

// StartWithBelt implements tracer.Tracer.
func (t *TracerImpl) StartWithBelt(belt *belt.Belt, name string, options ...tracer.SpanOption) (tracer.Span, *belt.Belt) {
	return t.newSpanBelt(belt, name, nil, options...)

}

// StartChildWithBelt implements tracer.Tracer.
func (t *TracerImpl) StartChildWithBelt(belt *belt.Belt, name string, options ...tracer.SpanOption) (tracer.Span, *belt.Belt) {
	return t.newSpanBelt(belt, name, tracer.SpanFromBelt(belt), options...)
}

// StartWithCtx implements tracer.Tracer.
func (t *TracerImpl) StartWithCtx(ctx context.Context, name string, options ...tracer.SpanOption) (tracer.Span, context.Context) {
	return t.newSpanCtx(ctx, name, nil, options...)

}

// StartChildWithCtx implements tracer.Tracer.
func (t *TracerImpl) StartChildWithCtx(ctx context.Context, name string, options ...tracer.SpanOption) (tracer.Span, context.Context) {
	return t.newSpanCtx(ctx, name, tracer.SpanFromCtx(ctx), options...)
}

// Flush implements tracer.Tracer.
func (t *TracerImpl) Flush(ctx context.Context) {
	t.Logger.Flush(ctx)
}

func (t *TracerImpl) send(span *SpanImpl) {
	fields := span.FieldsValue.
		WithField(FieldNameName, span.NameValue, FieldPropertyName).
		WithField(FieldNameDuration, span.Duration, FieldPropertyDuration).
		WithField(FieldNameStartTS, span.StartTSValue, FieldPropertyStartTS)

	if span.ParentValue != nil {
		fields = fields.WithField(FieldNameParent, span.ParentValue, FieldPropertyParent)
	}

	if len(span.Events) > 0 {
		eventsAsFields := make(field.Fields, 0, len(span.Events))
		for _, ev := range span.Events {
			eventsAsFields = append(eventsAsFields, field.Field{
				Key:        FieldNamePrefixEvent + ev.Name,
				Value:      ev.Timestamp,
				Properties: field.Properties{FieldPropertyEvent},
			})
		}
		fields = fields.WithFields(eventsAsFields)
	}

	t.Logger.LogFields(t.LevelFunc(span), "timespan finished", fields)
}

type fieldPropertyEnum uint

const (
	// FieldPropertyStartTS is a field.Property which marks the field as a start timestamp.
	FieldPropertyStartTS = fieldPropertyEnum(iota + 1)

	// FieldPropertyDuration is a field.Property which marks the field as a duration.
	FieldPropertyDuration

	// FieldPropertyName is a field.Property which marks the field as the span name
	FieldPropertyName

	// FieldPropertyEvent is a field.Property which marks the field as the timestamp of an event.
	FieldPropertyEvent

	// FieldPropertyParent is a field.Property which marks the field as the parent span.
	FieldPropertyParent
)

// SpanFromLogEntry converts a logger Entry to a Span.
func SpanFromLogEntry(entry *loggertypes.Entry) *SpanImpl {
	var (
		fields   field.Fields
		startTS  time.Time
		duration time.Duration
		parent   tracer.Span
		name     string
		events   []Event
	)
	entry.Fields.ForEachField(func(f *field.Field) bool {
		if len(f.Properties) != 1 {
			fields = append(fields, *f)
			return true
		}
		switch f.Properties[0] {
		case FieldPropertyStartTS:
			startTS = f.Value.(time.Time)
		case FieldPropertyParent:
			parent = f.Value.(tracer.Span)
		case FieldPropertyDuration:
			duration = f.Value.(time.Duration)
		case FieldPropertyName:
			name = f.Value.(string)
		case FieldPropertyEvent:
			events = append(events, Event{
				Timestamp: f.Value.(time.Time),
				Name:      f.Key[len(FieldNamePrefixEvent):],
			})
		default:
			fields = append(fields, *f)
		}
		return true
	})
	span := &SpanImpl{
		StartTSValue: startTS,
		Duration:     duration,
		ParentValue:  parent,
		FieldsValue:  (*field.FieldsChain)(nil).WithFields(fields),
		NameValue:    name,
		Events:       events,
	}
	return span
}

var (
	// FieldNamePrefixEvent is the field name prefix for events reported through Annotate method.
	FieldNamePrefixEvent = "event_"

	// FieldNameName is the field name for the span name.
	FieldNameName = "name"

	// FieldNameDuration is the field name for the duration.
	FieldNameDuration = "duration"

	// FieldNameParent is the field name for the parent span.
	FieldNameParent = "parent"

	// FieldNameStartTS is the field name for the span start timestamp.
	FieldNameStartTS = "start_ts"
)

// Default is the overridable function which returns a logger-based Tracer
// with the default configuration.
var Default = func() tracer.Tracer {
	return New(logger.Default(), func(_ *SpanImpl) loggertypes.Level {
		return logger.LevelTrace
	})
}

// New returns a new instance of TracerImpl (a tracer.Tracer implementation based on a given logger.Logger).
func New(logger loggertypes.Logger, levelFunc func(span *SpanImpl) loggertypes.Level) *TracerImpl {
	return &TracerImpl{
		Logger:    logger,
		LevelFunc: levelFunc,
	}
}
