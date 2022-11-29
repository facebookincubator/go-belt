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

package stdlib

import (
	"log"

	"github.com/facebookincubator/go-belt/tool/logger/adapter"
	"github.com/facebookincubator/go-belt/tool/logger/types"
)

// DefaultLevel is an overridable default logging level for
// loggers returned by function Default.
var DefaultLevel = types.LevelWarning

// NewEmitter returns a Emitter given a logger of the standard package "log"
func NewEmitter(stdLogger *log.Logger, opts ...types.Option) types.Emitter {
	return adapter.EmitterFromPrintfer(stdLogger, opts...)
}

var (
	defaultLogger = New(log.Default(), DefaultLevel)
)

// Default returns the default Logger based on the standard package "log".
var Default = func() types.Logger {
	return defaultLogger
}

// New returns a Logger using given a logger of the standard package "log".
func New(stdLogger *log.Logger, level types.Level, opts ...types.Option) types.Logger {
	return adapter.LoggerFromPrintfer(stdLogger, opts...).WithLevel(level)
}
