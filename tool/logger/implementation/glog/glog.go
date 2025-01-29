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

package glog

import (
	"github.com/facebookincubator/go-belt/tool/logger"
	"github.com/facebookincubator/go-belt/tool/logger/adapter"
	"github.com/facebookincubator/go-belt/tool/logger/types"
	"github.com/golang/glog"
)

// GLog is a types.Emitter implementation based on glog.
type GLog struct{}

var _ types.Emitter = (*GLog)(nil)

// NewEmitter returns a new instance of a types.Emitter based
// on the global instance of glog (the upstream glog does not have
// a non-global instance, unfortunately).
func NewEmitter(_ ...types.Option) GLog {
	return GLog{}
}

// New returns a new instance a types.Logger based on glog.
func New(opts ...types.Option) types.Logger {
	return adapter.LoggerFromEmitter(NewEmitter(opts...), opts...)
}

// Flush implements types.Emitter
func (l GLog) Flush() {
	glog.Flush()
}

// Emit implements types.Emitter
func (l GLog) Emit(entry *types.Entry) {
	switch entry.Level {
	case logger.LevelTrace, logger.LevelDebug, logger.LevelInfo:
		glog.Info(entry.Message)
	case logger.LevelWarning:
		glog.Warning(entry.Message)
	case logger.LevelError:
		glog.Error(entry.Message)
	case logger.LevelPanic:
		glog.Fatal(entry.Message)
	case logger.LevelFatal:
		glog.Exit(entry.Message)
	default:
		glog.Info("[UNKNOWN LOGGING LEVEL] " + entry.Message)
	}
}

// Severity is the internal glog's logging level.
type Severity glog.Level

// SeverityFromLevel converts logger.Level to Severity.
func SeverityFromLevel(level types.Level) Severity {
	switch level {
	case types.LevelTrace, types.LevelDebug:
		return 0
	case types.LevelInfo:
		return 0
	case types.LevelWarning:
		return 1
	case types.LevelError:
		return 2
	case types.LevelPanic, types.LevelFatal:
		return 3
	default:
		return 0
	}
}
