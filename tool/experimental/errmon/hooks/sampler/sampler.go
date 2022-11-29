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
	"github.com/facebookincubator/go-belt"
	"github.com/facebookincubator/go-belt/pkg/field"
	"github.com/facebookincubator/go-belt/pkg/sampler"
	errmontypes "github.com/facebookincubator/go-belt/tool/experimental/errmon/types"
)

var (
	// DefaultIsSampledKey is the default field name to store the information if the entry was sampled out
	DefaultIsSampledKey = "is_sampled"
)

// Sampled is true when the entry was sampled out (and false when it was not).
type Sampled bool

// Sampler is a PreHook for an ErrorMonitor to sample the issued events.
type Sampler struct {
	sampler.Sampler
	IsSampledKey field.Key
}

// New returns a new instance of Sampler.
func New(sampler sampler.Sampler, isSampledKey field.Key) *Sampler {
	return &Sampler{
		Sampler:      sampler,
		IsSampledKey: isSampledKey,
	}
}

var _ errmontypes.PreHook = (*Sampler)(nil)

// ProcessInputError implements errmon.PreHook.
func (hook *Sampler) ProcessInputError(traceIDs belt.TraceIDs, _ error) errmontypes.PreHookResult {
	return hook.processInput(traceIDs)
}

// ProcessInputPanic implements errmon.PreHook.
func (hook *Sampler) ProcessInputPanic(traceIDs belt.TraceIDs, _ any) errmontypes.PreHookResult {
	return hook.processInput(traceIDs)
}

func (hook *Sampler) processInput(traceIDs belt.TraceIDs) errmontypes.PreHookResult {
	if hook.ShouldStay(traceIDs) {
		return errmontypes.PreHookResult{
			ExtraFields: &field.Field{Key: hook.isSampledKey(), Value: Sampled(false)},
		}
	}
	return errmontypes.PreHookResult{
		Skip: true,
	}
}

func (hook *Sampler) isSampledKey() field.Key {
	if hook.IsSampledKey != "" {
		return hook.IsSampledKey
	}
	return DefaultIsSampledKey
}
