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

package logrus

import (
	"runtime"

	"github.com/sirupsen/logrus"
)

type logrusCtxKeyCallerT struct{}

var LogrusCtxKeyCaller = logrusCtxKeyCallerT{}

// RestoreCallerHook is a logrus.Hook which sets the Caller
// from a Context.
type RestoreCallerHook struct{}

var _ logrus.Hook = (*RestoreCallerHook)(nil)

func newRestoreCallerHook() RestoreCallerHook {
	return RestoreCallerHook{}
}

func (RestoreCallerHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
func (RestoreCallerHook) Fire(entry *logrus.Entry) error {

	// restoring the Caller:
	caller := entry.Context.Value(LogrusCtxKeyCaller)
	if caller == nil {
		// no Caller was set
		return nil
	}
	entry.Caller, _ = caller.(*runtime.Frame)

	// working around the check in standard formatters:
	origLogger := entry.Logger
	entry.Logger = &logrus.Logger{
		Out:          origLogger.Out,
		Hooks:        origLogger.Hooks,
		Formatter:    origLogger.Formatter,
		ReportCaller: true,
		Level:        origLogger.Level,
		ExitFunc:     origLogger.ExitFunc,
		BufferPool:   origLogger.BufferPool,
	}

	return nil
}
