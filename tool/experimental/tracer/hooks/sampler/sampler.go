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

package sampler

import (
	"github.com/facebookincubator/go-belt/pkg/sampler"
	"github.com/facebookincubator/go-belt/tool/experimental/tracer"
)

// Hook is a types.Hook implementation which samples the spans.
// This is supposed to be used for performance reasons in high performance applications.
//
// Setup example:
//
//	import (
//		"github.com/facebookincubator/go-belt/pkg/sampler"
//		samplerhook "github.com/facebookincubator/go-belt/tool/experimental/tracer/hooks/sampler"
//		"github.com/facebookincubator/go-belt/tool/experimental/tracer/implementation/zipkin"
//	)
//
//	func main() {
//		...
//		ctx = tracer.CtxWithTracer(ctx, zipkin.New(zipkinClient).WithPreHooks(
//			samplerhook.NewSamplerPreHook(sampler.RandomSampler(0.1)),
//		))
//		...
//	}
type Hook struct {
	sampler.Sampler
}

var _ tracer.Hook = (*Hook)(nil)

// ProcessSpan implements tracer.Hook
func (hook *Hook) ProcessSpan(span tracer.Span) bool {
	if span.Parent() != nil {
		// already decided on the parent, that we will log these spans
		return true
	}
	return hook.Sampler.ShouldStay(span.TraceIDs())
}

// Flush implements tracer.Hook
func (hook *Hook) Flush() {}
