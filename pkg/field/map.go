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

package field

import (
	"fmt"
)

// Map is just a handy low performance way to pass fields
// similar to how it is done in logrus.
type Map[V any] map[Key]V

var _ AbstractFields = (Map[Value])(nil)

// Len implements AbstractFields.
func (m Map[V]) Len() int {
	return len(m)
}

// ForEachField implements AbstractFields.
func (m Map[V]) ForEachField(callback func(f *Field) bool) bool {
	var f Field
	for k, v := range m {
		f.Key = k
		f.Value = v
		if !callback(&f) {
			return false
		}
	}
	return true
}

// Map is an abstract map which implements field.AbstractFields.
//
// It differs from Map, because it also accepts non-string keys.
type MapGeneric[K comparable, V any] map[K]V

var _ AbstractFields = (MapGeneric[string, any])(nil)

// ForEachField implements field.AbstractFields.
func (m MapGeneric[K, V]) ForEachField(callback func(f *Field) bool) bool {
	var f Field
	for k, v := range m {
		f.Key = fmt.Sprint(k)
		f.Value = v
		if !callback(&f) {
			return false
		}
	}
	return true
}

// Len implements field.AbstractFields.
func (m MapGeneric[K, V]) Len() int {
	return len(m)
}
