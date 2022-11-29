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
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/facebookincubator/go-belt/pkg/field"
	"github.com/facebookincubator/go-belt/tool/logger/types"
	"github.com/sirupsen/logrus"
)

func BenchmarkWithFields(b *testing.B) {
	logrusLogger := logrus.New()
	logrusLogger.Out = io.Discard
	logrusLogger.Level = logrus.TraceLevel

	for depth := 1; depth <= 256; depth *= 2 {
		var keys [256][4]string
		for i, _keys := range keys {
			for j := range _keys {
				keys[i][j] = fmt.Sprintf("key %d:%d", i, j)
			}
		}
		b.Run(fmt.Sprintf("depth%d", depth), func(b *testing.B) {
			for _, callLog := range []bool{false, true} {
				b.Run(fmt.Sprintf("callLog-%v", callLog), func(b *testing.B) {
					b.Run("bare_logrus", func(b *testing.B) {
						lOrig := logrusLogger.WithContext(context.Background())
						b.ReportAllocs()
						b.ResetTimer()
						for i := 0; i < b.N; i++ {
							l := lOrig
							for num := 0; num < depth; num++ {
								l = l.WithField(keys[num][0], "some value")
								l = l.WithFields(logrus.Fields{
									keys[num][1]: "more values 1",
									keys[num][2]: "more values 2",
									keys[num][3]: 3,
								})
							}
							if callLog {
								l.Logf(logrus.ErrorLevel, "unit-test")
							}
						}
					})
					b.Run("adapted_logrus", func(b *testing.B) {
						lOrig := New(logrusLogger)
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
								l.Logf(types.LevelError, "unit-test")
							}
						}
					})
				})
			}
		})
	}
}
