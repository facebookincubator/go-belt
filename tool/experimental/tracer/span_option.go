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

	"github.com/facebookincubator/go-belt/pkg/field"
)

// SpanOption is an option which modifies a Span.
//
// If some option is not supported by a specific
// implementation of the Tracer, then it is just ignored.
type SpanOption interface {
	_x(*Span)
}

// SpanOptionRole defines who is the issuer of this specific Span,
// which might be useful in some implementations of distributed tracing.
type SpanOptionRole string

const (
	spanOptionRoleUndefined = SpanOptionRole("") //nolint:deadcode,unused,varcheck
	RoleClient              = SpanOptionRole("client")
	RoleServer              = SpanOptionRole("server")
)

func (SpanOptionRole) _x(*Span) {}

// SpanOptionStart overrides the starting timestamp of the Span.
type SpanOptionStart time.Time

func (SpanOptionStart) _x(*Span) {}

// SpanOptionAddFields adds more structured fields to the Span.
//
// Note: some distributing tracers may not support logging this data.
type SpanOptionAddFields field.Fields

func (SpanOptionAddFields) _x(*Span) {}

// SpanOptionResetFields resets structured fields within the Span to an empty set.
type SpanOptionResetFields struct{}

func (SpanOptionResetFields) _x(*Span) {}
