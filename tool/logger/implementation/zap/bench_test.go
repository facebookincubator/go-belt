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
	"os"
	"testing"

	"github.com/facebookincubator/go-belt/pkg/field"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapConfig struct {
	Name   string
	Logger *zap.Logger
}

func zapConfigs(t interface{ Fatal(...any) }) []zapConfig {
	var err error

	stdErr := os.Stderr
	os.Stderr, err = os.Open("/dev/null")
	if err != nil {
		t.Fatal(err)
	}
	zapProd, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}
	zapDev, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}
	os.Stderr = stdErr

	return []zapConfig{
		{
			Name:   "prod",
			Logger: zapProd,
		},
		{
			Name:   "dev",
			Logger: zapDev,
		},
		{
			Name:   "nop",
			Logger: zap.NewNop(),
		},
	}
}

func Benchmark(b *testing.B) {
	for _, zapConfig := range zapConfigs(b) {
		b.Run(zapConfig.Name, func(b *testing.B) {

			for depth := 0; depth <= 256; depth = 1 + depth*4/3 {
				var keys [256][4]string
				for i, _keys := range keys {
					for j := range _keys {
						keys[i][j] = fmt.Sprintf("key %d:%d", i, j)
					}
				}
				b.Run(fmt.Sprintf("depth%d", depth), func(b *testing.B) {
					b.Run("WithField", func(b *testing.B) {
						for _, callLog := range []bool{false, true} {
							b.Run(fmt.Sprintf("callLog-%v", callLog), func(b *testing.B) {
								b.Run("bare_zap", func(b *testing.B) {
									lOrig := zapConfig.Logger
									b.ReportAllocs()
									b.ResetTimer()
									for i := 0; i < b.N; i++ {
										l := lOrig
										for num := 0; num < depth; num++ {
											l = l.With(zap.Field{
												Key:    keys[num][0],
												Type:   zapcore.StringType,
												String: "some value",
											})
											l = l.With(
												zap.Field{
													Key:    keys[num][1],
													Type:   zapcore.StringType,
													String: "more values 1",
												},
												zap.Field{
													Key:    keys[num][2],
													Type:   zapcore.StringType,
													String: "more values 2",
												},
												zap.Field{
													Key:     keys[num][3],
													Type:    zapcore.Int64Type,
													Integer: 3,
												},
											)
										}
										if callLog {
											l.Error("unit-test")
										}
									}
								})
								b.Run("adapted_zap", func(b *testing.B) {
									lOrig := New(zapConfig.Logger)
									b.ReportAllocs()
									b.ResetTimer()
									for i := 0; i < b.N; i++ {
										l := lOrig
										for num := 0; num < depth; num++ {
											l = l.WithField(keys[num][0], "some value")
											l = l.WithFields(field.Fields{
												{Key: keys[num][1], Value: "more values 1"},
												{Key: keys[num][2], Value: "more values 2"},
												{Key: keys[num][3], Value: 3},
											})
										}
										if callLog {
											l.Errorf("unit-test")
										}
									}
								})
							})
						}
					})
					b.Run("Log", func(b *testing.B) {
						zapLogger := zapConfig.Logger
						for num := 0; num < depth; num++ {
							zapLogger = zapLogger.With(zap.Field{
								Key:    keys[num][0],
								Type:   zapcore.StringType,
								String: "some value",
							})
						}

						for fieldNum := 0; fieldNum < 8; fieldNum++ {
							b.Run(fmt.Sprintf("fields-%d", fieldNum), func(b *testing.B) {
								b.Run("bare_zap", func(b *testing.B) {
									l := zapLogger
									var fields []zap.Field
									for i := 0; i < fieldNum; i++ {
										fields = append(fields, zap.Field{
											Key:    fmt.Sprintf("key %d", i),
											Type:   zapcore.StringType,
											String: fmt.Sprintf("value %d", i),
										})
									}
									b.ReportAllocs()
									b.ResetTimer()
									for i := 0; i < b.N; i++ {
										l.Error("unit-test", fields...)
									}
								})
								b.Run("adapted_zap", func(b *testing.B) {
									l := New(zapLogger)
									var fields field.Fields
									for i := 0; i < fieldNum; i++ {
										fields = append(fields, field.Field{
											Key:   fmt.Sprintf("key %d", i),
											Value: fmt.Sprintf("value %d", i),
										})
									}
									b.ReportAllocs()
									b.ResetTimer()
									for i := 0; i < b.N; i++ {
										l.ErrorFields("unit-test", &fields)
									}
								})
							})
						}
					})
				})
			}
		})
	}
}
