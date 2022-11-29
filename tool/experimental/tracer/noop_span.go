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

package tracer

import (
	"time"

	"github.com/facebookincubator/go-belt"
	"github.com/facebookincubator/go-belt/pkg/field"
)

// NoopSpan is a no-op implementation of a Span. Supposed
// to be used for sampled out spans.
type NoopSpan struct {
	ParentValue  Span
	NameValue    string
	StartTSValue time.Time
}

var _ Span = (*NoopSpan)(nil)

// NewNoopSpan returns a new instance of a NoopSpan.
func NewNoopSpan(name string, parent Span, now time.Time) *NoopSpan {
	return &NoopSpan{
		ParentValue:  parent,
		NameValue:    name,
		StartTSValue: now,
	}
}

// ID implements Span.
func (*NoopSpan) ID() any {
	return nil
}

// TraceIDs implements Span.
func (*NoopSpan) TraceIDs() belt.TraceIDs {
	return nil
}

// Fields implements Span.
func (*NoopSpan) Fields() field.AbstractFields {
	return nil
}

// Parent implements Span.
func (s *NoopSpan) Parent() Span {
	return s.ParentValue
}

// SetName implements Span.
func (s *NoopSpan) SetName(name string) { s.NameValue = name }

// Name implements Span.
func (s *NoopSpan) Name() string { return s.NameValue }

// StartTS implements Span.
func (s *NoopSpan) StartTS() time.Time { return s.StartTSValue }

// Annotate implements Span.
func (*NoopSpan) Annotate(time.Time, string) {}

// SetField implements Span.
func (*NoopSpan) SetField(field.Key, field.Value) {}

// SetFields implements Span.
func (*NoopSpan) SetFields(field.AbstractFields) {}

// Finish implements Span.
func (*NoopSpan) Finish() {}

// FinishWithDuration implements Span.
func (*NoopSpan) FinishWithDuration(time.Duration) {}

// Flush implements Span.
func (*NoopSpan) Flush() {}

// IsNoopSpan returns true if the given Span is a NoopSpan.
func IsNoopSpan(span Span) bool {
	_, isNoop := span.(*NoopSpan)
	return isNoop
}
