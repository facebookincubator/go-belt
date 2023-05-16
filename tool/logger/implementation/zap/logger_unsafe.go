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
	"math"
	"sync/atomic"
	"unsafe"

	"github.com/facebookincubator/go-belt/pkg/field"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func fieldsToZap(length int, forEachField func(callback func(f *field.Field) bool) bool) *[]zap.Field {
	if length == 0 {
		return nil
	}
	result := acquireZapFields()
	if cap(*result) < length {
		*result = make([]zap.Field, 0, length)
	}
	if len(*result) != 0 {
		panic(len(*result))
	}
	count := 0
	appendToResult := func(f *field.Field) bool {
		count++
		if count > length {
			return false
		}
		zapField := zap.Field{
			Key: f.Key,
		}
		switch value := f.Value.(type) {
		case string:
			zapField.Type = zapcore.StringType
			zapField.String = value
		case int:
			zapField.Type = zapcore.Int64Type
			zapField.Integer = int64(value)
		case float64:
			zapField.Type = zapcore.Float64Type
			zapField.Integer = int64(math.Float64bits(value))
		case error:
			zapField.Type = zapcore.ErrorType
			zapField.Interface = f.Value
		default:
			zapField.Type = zapcore.ReflectType
			zapField.Interface = f.Value
		}
		*result = append(*result, zapField)
		return true
	}

	// See the assumption described in the "WARNING" message of field.ForEachFieldser.
	forEachField(*(*func(f *field.Field) bool)(noescape(unsafe.Pointer(&appendToResult))))

	if len(*result) > length {
		panic("should not happen")
	}
	return result
}

func (l *Emitter) setZapLogger(logger *zap.Logger) {
	atomic.StorePointer((*unsafe.Pointer)((unsafe.Pointer)(&l.ZapLogger)), (unsafe.Pointer)(logger))
}

func (l *Emitter) getZapLogger() *zap.Logger {
	return (*zap.Logger)(atomic.LoadPointer((*unsafe.Pointer)(noescape((unsafe.Pointer)(&l.ZapLogger)))))
}
