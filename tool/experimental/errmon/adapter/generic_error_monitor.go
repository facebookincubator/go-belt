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

package adapter

import (
	"time"

	"github.com/facebookincubator/go-belt"
	"github.com/facebookincubator/go-belt/pkg/field"
	"github.com/facebookincubator/go-belt/pkg/runtime"
	errmontypes "github.com/facebookincubator/go-belt/tool/experimental/errmon/types"
	"github.com/facebookincubator/go-belt/tool/experimental/tracer"
	loggertypes "github.com/facebookincubator/go-belt/tool/logger/types"
	"github.com/go-ng/slices"
)

// ErrorMonitorFromEmitter wraps a Emitter with a generic implementation of of ErrorMonitor.
func ErrorMonitorFromEmitter(
	emitter Emitter,
	callerFrameFilter runtime.PCFilter,
) ErrorMonitor {
	return &GenericErrorMonitor{
		EmitterValue:      emitter,
		CallerFrameFilter: callerFrameFilter,
	}
}

// GenericErrorMonitor implements ErrorMonitor given a Emitter.
type GenericErrorMonitor struct {
	EmitterValue      Emitter
	CallerFrameFilter runtime.PCFilter
	ContextFields     *FieldsChain
	TraceIDs          TraceIDs
	PreHooks          PreHooks
	Hooks             Hooks
}

var _ ErrorMonitor = (*GenericErrorMonitor)(nil)

func (h GenericErrorMonitor) clone() *GenericErrorMonitor {
	return &h
}

func (h *GenericErrorMonitor) newEvent(belt *belt.Belt, extraFields field.AbstractFields) *Event {
	fields := h.ContextFields
	if extraFields != nil {
		fields = fields.WithFields(extraFields)
	}
	ev := &Event{
		Entry: loggertypes.Entry{
			Timestamp: time.Now(),
			Fields:    fields,
			TraceIDs:  h.TraceIDs,
		},
		ID: errmontypes.RandomEventID(),
	}

	// stack trace
	ev.Exception.StackTrace = runtime.CallerStackTrace(h.CallerFrameFilter)

	// caller
	var pcs [1]runtime.PC
	if ev.StackTrace.ProgramCounters(pcs[:]) != 0 {
		ev.Entry.Caller = pcs[0]
	}

	// goroutines
	ev.Goroutines, ev.CurrentGoroutineID = getGoroutines()

	// spans
	span := tracer.SpanFromBelt(belt)
	for span != nil {
		ev.Spans = append(ev.Spans, span)
		span = span.Parent()
	}
	slices.Reverse(ev.Spans)

	return ev
}

// Emitter implements errmon.ErrorMonitor
func (h *GenericErrorMonitor) Emitter() Emitter {
	return h.EmitterValue
}

// WithPreHooks implements errmon.ErrorMonitor
func (h *GenericErrorMonitor) WithPreHooks(preHooks ...PreHook) ErrorMonitor {
	clone := h.clone()
	clone.PreHooks = preHooks
	return clone
}

// WithHooks implements errmon.ErrorMonitor
func (h *GenericErrorMonitor) WithHooks(hooks ...Hook) ErrorMonitor {
	clone := h.clone()
	clone.Hooks = hooks
	return clone
}

// ObserveError implements errmon.ErrorMonitor
func (h *GenericErrorMonitor) ObserveError(belt *belt.Belt, err error) *Event {
	if err == nil {
		return nil
	}
	preHooksResult := h.PreHooks.ProcessInputError(h.TraceIDs, err)
	if preHooksResult.Skip {
		return nil
	}
	ev := h.newEvent(belt, preHooksResult.ExtraFields)
	if ev == nil {
		return nil
	}
	ev.Exception.Error = err
	ev.Entry.Level = loggertypes.LevelError
	ev.Entry.Properties = []loggertypes.EntryProperty{errmontypes.EntryPropertyErrorMonitoringEventEntry, errmontypes.EntryPropertyErrorEvent}
	if !h.Hooks.Process(ev) {
		return nil
	}
	h.EmitterValue.Emit(ev)
	return ev
}

// ObserveRecover implements errmon.ErrorMonitor
func (h *GenericErrorMonitor) ObserveRecover(belt *belt.Belt, recoverResult any) *Event {
	if recoverResult == nil {
		// no panic
		return nil
	}
	preHooksResult := h.PreHooks.ProcessInputPanic(h.TraceIDs, recoverResult)
	if preHooksResult.Skip {
		return nil
	}
	ev := h.newEvent(belt, preHooksResult.ExtraFields)
	if ev == nil {
		return nil
	}
	ev.Exception.IsPanic = true
	ev.Exception.PanicValue = recoverResult
	ev.Entry.Level = loggertypes.LevelPanic
	ev.Entry.Properties = []loggertypes.EntryProperty{errmontypes.EntryPropertyErrorMonitoringEventEntry, errmontypes.EntryPropertyPanicEvent}
	if !h.Hooks.Process(ev) {
		return nil
	}
	h.EmitterValue.Emit(ev)
	return ev
}

// WithContextFields implements errmon.ErrorMonitor
func (h *GenericErrorMonitor) WithContextFields(allFields *FieldsChain, _ int) belt.Tool {
	clone := h.clone()
	clone.ContextFields = allFields
	return clone
}

// WithTraceIDs implements errmon.ErrorMonitor
func (h *GenericErrorMonitor) WithTraceIDs(allTraceIDs TraceIDs, _ int) belt.Tool {
	clone := h.clone()
	clone.TraceIDs = allTraceIDs
	return clone
}

// Flush implements metrics.Metrics (or more specifically belt.Tool).
func (*GenericErrorMonitor) Flush() {}
