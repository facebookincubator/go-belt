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

package zap

import (
	"fmt"

	"github.com/facebookincubator/go-belt/tool/logger/types"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LevelToZap converts logger.Level to zap's logging level.
func LevelToZap(level types.Level) zapcore.Level {
	switch level {
	case types.LevelTrace:
		return zap.DebugLevel
	case types.LevelDebug:
		return zap.DebugLevel
	case types.LevelInfo:
		return zap.InfoLevel
	case types.LevelWarning:
		return zap.WarnLevel
	case types.LevelError:
		return zap.ErrorLevel
	case types.LevelPanic:
		return zap.PanicLevel
	case types.LevelFatal:
		return zap.FatalLevel
	}
	panic(fmt.Errorf("unexpected level: %v", level))
}

// LevelFromZap converts zap's logging level to logger.Level.
func LevelFromZap(level zapcore.Level) types.Level {
	switch level {
	case zap.DebugLevel - 1:
		return types.LevelTrace
	case zap.DebugLevel:
		return types.LevelDebug
	case zap.InfoLevel:
		return types.LevelInfo
	case zap.WarnLevel:
		return types.LevelWarning
	case zap.ErrorLevel:
		return types.LevelError
	case zap.PanicLevel:
		return types.LevelPanic
	case zap.FatalLevel:
		return types.LevelFatal
	}
	panic(fmt.Errorf("unexpected level: %v", level))
}
