// Copyright 2023 Meta Platforms, Inc. and affiliates.
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

package slog

import (
	"github.com/facebookincubator/go-belt/tool/logger/types"
	"golang.org/x/exp/slog"
)

// LevelToSlog converts logger.Level to slog.Level.
func LevelToSlog(level types.Level) slog.Level {
	switch level {
	case types.LevelTrace:
		return slog.LevelDebug - 1
	case types.LevelDebug:
		return slog.LevelDebug
	case types.LevelInfo:
		return slog.LevelInfo
	case types.LevelWarning:
		return slog.LevelWarn
	case types.LevelError:
		return slog.LevelError
	case types.LevelPanic:
		return slog.LevelError + 1
	case types.LevelFatal:
		return slog.LevelError + 2
	}

	// not mapped logging level is an error per se:
	return slog.LevelError
}

// LevelFromSlog converts slog.Level to logger.Level.
func LevelFromSlog(level slog.Level) types.Level {
	switch level {
	case slog.LevelDebug - 1:
		return types.LevelTrace
	case slog.LevelDebug:
		return types.LevelDebug
	case slog.LevelInfo:
		return types.LevelInfo
	case slog.LevelWarn:
		return types.LevelWarning
	case slog.LevelError:
		return types.LevelError
	case slog.LevelError + 1:
		return types.LevelPanic
	case slog.LevelError + 2:
		return types.LevelFatal
	}

	// not mapped logging level is an error per se:
	return types.LevelError
}
