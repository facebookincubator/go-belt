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

// Copyright (c) Facebook, Inc. and its affiliates.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package field

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testContextField = &FieldsChain{
		parent: &FieldsChain{
			parent: &FieldsChain{
				oneField: Field{
					Key:   "5",
					Value: 6,
				},
			},
			multipleFields: Fields{{
				Key:   "3",
				Value: 4,
			}},
		},
		oneField: Field{
			Key:   "1",
			Value: 2,
		},
	}

	_v *FieldsChain
)

func dummyFields(count uint) Fields {
	result := make(Fields, count)
	for idx := range result {
		result[idx] = Field{Key: fmt.Sprintf("%d", idx), Value: idx}
	}
	return result
}

func BenchmarkFields_CopyAndAddOne(b *testing.B) {
	fields := Gather(testContextField)
	for _, cloneDepth := range []uint{1, 2, 4, 8, 16, 32, 64} {
		_fields := dummyFields(cloneDepth)
		b.Run(fmt.Sprintf("cloneDepth%d", cloneDepth), func(b *testing.B) {
			initialFields := fields.Copy()
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				f := initialFields
				for i := 0; i < int(cloneDepth); i++ {
					f = f.Copy()
					f = append(f, _fields[i])
				}
			}
		})
	}
}

func BenchmarkFieldsChain_CloneAndAddOneAnd(b *testing.B) {
	// Last benchmark: BenchmarkFieldsChain_CloneAndAddOneAnd/cloneDepth1/withGather-false-16  40.44 ns/op  96 B/op  1 allocs/op
	//
	// This allocation and CPU consumption could be reduced via sync.Pool,
	// but for now it is a premature optimization (there was no request on
	// such performance).
	for _, cloneDepth := range []uint{1, 2, 4, 8, 16, 32, 64} {
		keys := make([]string, cloneDepth)
		for idx := range keys {
			keys[idx] = fmt.Sprintf("%d", idx)
		}
		b.Run(fmt.Sprintf("cloneDepth%d", cloneDepth), func(b *testing.B) {
			for _, withGather := range []bool{false, true} {
				b.Run(fmt.Sprintf("withGather-%v", withGather), func(b *testing.B) {
					b.ReportAllocs()
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						f := testContextField
						for i := 0; i < int(cloneDepth); i++ {
							f = f.WithField(keys[i], i)
						}
						if withGather {
							_ = Gather(f)
						}
					}
				})
			}
		})
	}
}

func BenchmarkFieldsChain_CloneAndAddOneAsMultiple(b *testing.B) {
	fields := Fields{{
		Key:   "1",
		Value: "2",
	}}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_v = testContextField.WithFields(fields)
	}
}

func BenchmarkFieldsChain_CloneAndAddXAsMap(b *testing.B) {
	for x := 1; x < 32; x++ {
		b.Run(fmt.Sprintf("%d", x), func(b *testing.B) {
			m := map[Key]interface{}{}
			for i := 0; i < x; i++ {
				m[fmt.Sprintf("%d", i)] = i
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_v = testContextField.WithMap(m)
			}
		})
	}
}

func BenchmarkFieldsChain_Gather(b *testing.B) {
	// Last benchmark: BenchmarkFieldsChain_Compile-8  4787438  251 ns/op  336 B/op  2 allocs
	//
	// Could be performed the same optimization as for BenchmarkFieldsChain_CloneAndAddOne
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Gather(testContextField)
	}
}

func TestFieldsChain_Compile(t *testing.T) {
	require.Equal(t, Fields{
		{
			Key:   "1",
			Value: 2,
		},
		Field{
			Key:   "3",
			Value: 4,
		},
		Field{
			Key:   "5",
			Value: 6,
		},
	}, Gather(testContextField))
}
