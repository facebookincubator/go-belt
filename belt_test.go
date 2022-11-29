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

package belt

import (
	"fmt"
	"testing"
)

func BenchmarkBeltTraceIDs(b *testing.B) {
	for _, isSequential := range []bool{false, true} {
		b.Run(fmt.Sprintf("seq-%v", isSequential), func(b *testing.B) {
			for _, withTraceIDsCall := range []bool{false, true} {
				b.Run(fmt.Sprintf("compile-%v", withTraceIDsCall), func(b *testing.B) {
					for amount := 1; amount < 1024; amount *= 2 {
						b.Run(fmt.Sprintf("%d", amount), func(b *testing.B) {
							traceIDs := make([]TraceID, amount)
							for idx := range traceIDs {
								traceIDs[idx] = TraceID(fmt.Sprintf("%d", idx))
							}
							belt := New()
							b.ReportAllocs()
							b.ResetTimer()
							if isSequential {
								for i := 0; i < b.N; i++ {
									belt := belt.clone()
									for _, traceID := range traceIDs {
										belt = belt.WithTraceID(traceID)
									}
									if withTraceIDsCall {
										belt.TraceIDs()
									}
								}
							} else {
								for i := 0; i < b.N; i++ {
									clone := belt.WithTraceID(traceIDs...)
									if withTraceIDsCall {
										clone.TraceIDs()
									}
								}
							}
						})
					}
				})
			}
		})
	}
}
