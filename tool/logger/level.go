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
	"github.com/facebookincubator/go-belt/tool/logger/types"
)

// Level is used to define severity of messages to be reported.
type Level = types.Level

const (
	// LevelUndefined is the erroneous value of log-level which corresponds
	// to zero-value.
	LevelUndefined = types.LevelUndefined

	// LevelFatal will report about Fatalf-s only.
	LevelFatal = types.LevelFatal

	// LevelPanic will report about Panicf-s and Fatalf-s only.
	LevelPanic = types.LevelPanic

	// LevelError will report about Errorf-s, Panicf-s, ...
	LevelError = types.LevelError

	// LevelWarning will report about Warningf-s, Errorf-s, ...
	LevelWarning = types.LevelWarning

	// LevelInfo will report about Infof-s, Warningf-s, ...
	LevelInfo = types.LevelInfo

	// LevelDebug will report about Debugf-s, Infof-s, ...
	LevelDebug = types.LevelDebug

	// LevelTrace will report about Tracef-s, Debugf-s, ...
	LevelTrace = types.LevelTrace
)
