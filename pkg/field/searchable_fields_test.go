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

//go:build go_obs_experimental
// +build go_obs_experimental

package field

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchableFieldsDeduplicate(t *testing.T) {
	type testCase struct {
		Input  Fields
		Expect Fields
	}

	for testCaseID, testCase := range []testCase{
		{
			Input:  Fields{},
			Expect: Fields{},
		},
		{
			Input: Fields{
				{
					Key:   "1",
					Value: 2,
				},
			},
			Expect: Fields{
				{
					Key:   "1",
					Value: 2,
				},
			},
		},
		{
			Input: Fields{
				{
					Key:   "1",
					Value: 2,
				},
				{
					Key:   "1",
					Value: 3,
				},
			},
			Expect: Fields{
				{
					Key:   "1",
					Value: 3,
				},
			},
		},
		{
			Input: Fields{
				{
					Key:   "0",
					Value: 0,
				},
				{
					Key:   "1",
					Value: 2,
				},
				{
					Key:   "1",
					Value: 3,
				},
			},
			Expect: Fields{
				{
					Key:   "0",
					Value: 0,
				},
				{
					Key:   "1",
					Value: 3,
				},
			},
		},
		{
			Input: Fields{
				{
					Key:   "1",
					Value: 2,
				},
				{
					Key:   "1",
					Value: 3,
				},
				{
					Key:   "4",
					Value: 5,
				},
			},
			Expect: Fields{
				{
					Key:   "1",
					Value: 3,
				},
				{
					Key:   "4",
					Value: 5,
				},
			},
		},
		{
			Input: Fields{
				{
					Key:   "4",
					Value: 5,
				},
				{
					Key:   "1",
					Value: 2,
				},
				{
					Key:   "1",
					Value: 3,
				},
			},
			Expect: Fields{
				{
					Key:   "1",
					Value: 3,
				},
				{
					Key:   "4",
					Value: 5,
				},
			},
		},
	} {
		s := NewSearchableFields(testCase.Input)
		s.Build()
		assert.Equal(t, testCase.Expect, s.Fields(), fmt.Sprintf("test case #%d", testCaseID))
	}
}

func BenchmarkSearchableFieldsDeduplicateKeys(b *testing.B) {
	for _, isSorted := range []bool{false, true} {
		b.Run(fmt.Sprintf("isSorted-%v", isSorted), func(b *testing.B) {
			maxSize := 1024 * 1024
			for totalSize := 1; totalSize <= maxSize; totalSize *= 2 {
				b.Run(fmt.Sprintf("fields%d", totalSize), func(b *testing.B) {
					for duplicatesPlus1 := 1; duplicatesPlus1 <= maxSize; duplicatesPlus1 *= 2 {
						if duplicatesPlus1 > totalSize {
							b.Skip()
							return
						}
						duplicates := duplicatesPlus1 - 1
						b.Run(fmt.Sprintf("dups%d", duplicates), func(b *testing.B) {
							var fields SearchableFields
							fields.fields = make(Fields, totalSize)
							idxs := rand.Perm(totalSize)
							for idx := range idxs[duplicates:] {
								fields.fields[idx] = Field{
									Key:   fmt.Sprint(idx),
									Value: idx,
								}
							}
							for idx := range idxs[:duplicates] {
								dupOfIdx := rand.Intn(len(idxs[duplicates:]))
								fields.fields[idx] = fields.fields[idxs[dupOfIdx]]
							}

							if isSorted {
								fields.sort()
							}
							b.ReportAllocs()
							b.ResetTimer()
							for i := 0; i < b.N; i++ {
								b.StopTimer()
								fields := fields.Copy()
								b.StartTimer()
								fields.Build()
							}
						})
					}
				})
			}
		})
	}
}

func BenchmarkSearchableFieldsAddGet(b *testing.B) {
	for _, fieldsCount := range []uint{1, 2, 4, 8, 16, 32, 33, 64, 128, 256, 512, 1024} {
		b.Run(fmt.Sprintf("fieldCount%d", fieldsCount), func(b *testing.B) {
			for _, getCount := range []uint{0, 1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024} {
				b.Run(fmt.Sprintf("getCount%v", getCount), func(b *testing.B) {
					_fields := dummyFields(fieldsCount + getCount)
					b.ReportAllocs()
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						fields := NewSearchableFields(nil)
						for idx := uint(0); idx < fieldsCount; idx++ {
							fields.Add(_fields[idx])
						}
						for idx := 0; idx < int(getCount); idx++ {
							fields.Get(_fields[getCount].Key)
						}
					}
				})
			}
		})
	}
}

func BenchmarkSearchableFields_CopyAndAddOne(b *testing.B) {
	fields := Gather(testContextField)
	for _, withGet := range []bool{false, true} {
		b.Run(fmt.Sprintf("withGet-%v", withGet), func(b *testing.B) {
			for _, cloneDepth := range []uint{1, 2, 4, 8, 16, 32, 64} {
				_fields := dummyFields(cloneDepth)
				b.Run(fmt.Sprintf("cloneDepth%d", cloneDepth), func(b *testing.B) {
					initialFields := NewSearchableFields(fields.Copy())
					b.ReportAllocs()
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						f := initialFields
						for i := 0; i < int(cloneDepth); i++ {
							f = f.Copy()
							f.Add(_fields[i])
						}
						if withGet {
							f.Get(_fields[0].Key)
						}
					}
				})
			}
		})
	}
}
