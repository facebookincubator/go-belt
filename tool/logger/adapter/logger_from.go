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
	"github.com/facebookincubator/go-belt/tool/logger/types"
)

// LoggerFromAny makes the best effort to convert given logger of
// an arbitrary type to logger.Logger. Returns an untyped `nil` if
// was not successful.
func LoggerFromAny(anyLogger any, opts ...types.Option) types.Logger {
	if anyLogger == nil {
		return nil
	}

	switch l := anyLogger.(type) {
	case types.Logger:
		return l
	case CompactLogger:
		return LoggerFromCompactLogger(l, opts...)
	case types.Emitter:
		return LoggerFromEmitter(l, opts...)
	case Printfer:
		return LoggerFromPrintfer(l, opts...)
	case func(format string, args ...any):
		return LoggerFromPrintf(l, opts...)
	}

	return nil
}

// LoggerFromCompactLogger provides a generic logger.Logger based on a given types.CompactLogger.
func LoggerFromCompactLogger(compactLogger CompactLogger, opts ...types.Option) types.Logger {
	return GenericSugar{CompactLogger: compactLogger}
}

// LoggerFromEmitter provides a generic logger.Logger based on a given types.Emitter.
func LoggerFromEmitter(emitter types.Emitter, opts ...types.Option) types.Logger {
	cfg := types.Options(opts).Config()
	return LoggerFromCompactLogger(
		&GenericLogger{
			Emitters:      types.Emitters{emitter},
			GetCallerFunc: cfg.GetCallerFunc,
		},
		opts...,
	)
}

// EmitterFromPrintfer provides a naive types.Emitter based on a given Printfer.
func EmitterFromPrintfer(printfer Printfer, opts ...types.Option) types.Emitter {
	return PrintferEmitter{Printfer: printfer}
}

// EmitterFromPrintf provides a naive types.Emitter based on a given printf-function.
func EmitterFromPrintf(printf func(format string, args ...any), opts ...types.Option) types.Emitter {
	return EmitterFromPrintfer(printfWrap{Func: printf})
}

// LoggerFromPrintfer provides a generic naive logger.Logger based on a given Printfer.
func LoggerFromPrintfer(printfer Printfer, opts ...types.Option) types.Logger {
	return LoggerFromEmitter(EmitterFromPrintfer(printfer, opts...), opts...)
}

// LoggerFromPrintf provides a generic naive logger.Logger based on a given printf-function.
func LoggerFromPrintf(printf func(format string, args ...any), opts ...types.Option) types.Logger {
	return LoggerFromEmitter(EmitterFromPrintf(printf, opts...), opts...)
}
