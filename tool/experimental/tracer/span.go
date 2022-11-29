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

// Span is a single time interval.
type Span interface {
	// ID is the ID for the Span. It may have different types, depending
	// on specific implementation of the Tracer.
	//
	// Not all tracing systems may support that.
	ID() any

	// TraceIDs are the set if unique IDs associated with the current context.
	//
	// See the description of belt.TraceIDs.
	//
	// Not all tracing systems may support that.
	TraceIDs() belt.TraceIDs

	// Name is the name of the Span. Usually it explains what happens within
	// this span.
	Name() string

	// StartTS returns the timestamp of the beginning of the Span.
	//
	// Can be zero value if unknown. For example it could
	// happen due sampling, or it could just be not supported
	// byt a specific implementation.
	StartTS() time.Time

	// Fields is the set of structured data attached to the Span.
	Fields() field.AbstractFields

	// Parent is the parent Span. If there is no parent then an untyped nil is returned.
	Parent() Span

	// SetName sets the Name (see the description of method Name).
	SetName(string)

	// Annotate adds an event with the specified timestamp and description to the Span.
	//
	// Not all tracing systems may support that.
	Annotate(ts time.Time, description string)

	// SetField sets a structured field in the Span.
	//
	// Not all tracing systems may support that.
	SetField(field.Key, field.Value)

	// SetField set multiple structured fields at once in the Span.
	//
	// Not all tracing systems may support that.
	SetFields(field.AbstractFields)

	// Finish closes the Span using time.Now() as the ending timestamp and sends it.
	Finish()

	// FinishWithDuration closes the Span using startTS+duration as the ending timestamp and sends it.
	FinishWithDuration(duration time.Duration)

	// Flush forces to immediatelly empty all the buffers (send if something was delayed and so on).
	Flush()
}

var _ belt.Artifact = (Span)(nil)

// Spans is a collection of Span-s.
type Spans []Span

// Earliest returns the Span with the lowest start timestamp.
func (s Spans) Earliest() Span {
	if len(s) == 0 {
		return nil
	}
	earliest := s[0]
	for _, span := range s {
		if span.StartTS().Before(earliest.StartTS()) {
			earliest = span
		}
	}
	return earliest
}
