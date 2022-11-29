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
)

// Default is the (overridable) function which returns returns the default tracer.
//
// It is used by FromBelt and FromCtx functions if a Tracer is not set in the Belt.
var Default = func() Tracer {
	return noopTracer{}
}

// StartChildSpanFromBelt returns a child span given an Belt. The parent Span is extracted
// from the Belt. If one is not set, then a `nil` parent is set in the new Span.
//
// The Span and an Belt which includes this Span are returned.
func StartChildSpanFromBelt(belt *belt.Belt, name string, options ...SpanOption) (Span, *belt.Belt) {
	return FromBelt(belt).StartChildWithBelt(belt, name, options...)
}

// StartSpanFromBelt returns a span (without a parent) given an Belt.
//
// The Span and an Belt which includes this Span are returned.
func StartSpanFromBelt(belt *belt.Belt, name string, options ...SpanOption) (Span, *belt.Belt) {
	return FromBelt(belt).StartWithBelt(belt, name, options...)
}

// SpanFromBelt extracts a Span from a Belt (for example, previously derived using StartChildSpanFromBelt,
// StartChildSpanFromCtx or BeltWithSpan).
func SpanFromBelt(belt *belt.Belt) Span {
	loggerIface := belt.Artifacts().GetByID(ArtifactIDSpan)
	if loggerIface == nil {
		return NewNoopSpan("", nil, time.Time{})
	}
	return loggerIface.(Span)
}

// BeltWithSpan returns a Belt derivative with the specified Span.
func BeltWithSpan(belt *belt.Belt, span Span) *belt.Belt {
	return belt.WithArtifact(ArtifactIDSpan, span)
}

// FromBelt returns a Tracer given a Belt.
func FromBelt(belt *belt.Belt) Tracer {
	loggerIface := belt.Tools().GetByID(ToolID)
	if loggerIface == nil {
		return Default()
	}
	return loggerIface.(Tracer)
}

// BeltWithTracer returns a Belt derivative with the specified Tracer.
func BeltWithTracer(belt *belt.Belt, tracer Tracer) *belt.Belt {
	return belt.WithTool(ToolID, tracer)
}
