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
	"fmt"

	"github.com/facebookincubator/go-belt/pkg/field"
	errmonadapter "github.com/facebookincubator/go-belt/tool/experimental/errmon/adapter"
	errmontypes "github.com/facebookincubator/go-belt/tool/experimental/errmon/types"
	"github.com/facebookincubator/go-belt/tool/logger"
)

// Emitter is an implementation of errmon.Emitter which just prints the events
// through a given Logger.
type Emitter struct {
	Logger logger.Logger
}

// NewEmitter returns a new instance of Emitter.
func NewEmitter(logger logger.Logger) *Emitter {
	return &Emitter{
		Logger: logger,
	}
}

// New returns an instance of errmon.ErrorMonitor which prints all the events
// through a given Logger.
func New(logger logger.Logger, opts ...Option) errmontypes.ErrorMonitor {
	return errmonadapter.ErrorMonitorFromEmitter(
		NewEmitter(logger),
		options(opts).Config().CallerFrameFilter,
	)
}

// Default is an overridable function which returns an ErrorMonitor with the default configuration.
var Default = func() errmontypes.ErrorMonitor {
	return New(logger.Default())
}

// Flush implements errmon.Emitter
func (*Emitter) Flush() {}

// Emit implements errmon.Emitter
func (h *Emitter) Emit(ev *errmontypes.Event) {
	switch {
	case ev.Exception.IsPanic:
		ev.Entry.Message = fmt.Sprintf("got panic with argument: %v", ev.Exception.PanicValue)
	case ev.Exception.Error != nil:
		ev.Entry.Message = ev.Exception.Error.Error()
	}

	var stackTrace any
	if ev.Exception.StackTrace != nil {
		stackTrace = ev.Exception.StackTrace.String()
	}

	level := eventLevel(ev)
	ev.Entry.Level = level
	ev.Entry.Fields = field.Add(ev.Fields, field.Map[any]{
		"error_event.id":                    ev.ID,
		"error_event.external_ids":          ev.ExternalIDs,
		"error_event.exception.is_panic":    ev.Exception.IsPanic,
		"error_event.exception.panic_value": ev.Exception.PanicValue,
		"error_event.exception.error":       ev.Exception.Error,
		"error_event.exception.stack_trace": stackTrace,
		"error_event.spans":                 ev.Spans,
		"error_event.current_goroutine_id":  ev.CurrentGoroutineID,
		"error_event.goroutines":            ev.Goroutines,
	})
	h.Logger.Log(level, &ev.Entry)
}

func eventLevel(ev *errmontypes.Event) logger.Level {
	if ev.Level < logger.LevelError {
		return logger.LevelError
	}
	return ev.Level
}
