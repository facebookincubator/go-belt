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

// ForEachFieldser is an iterator through a collection of fields.
//
// Any object which implements this interface may be used as a source
// of fields. For example a generic parser of field values in a Logger
// checks if a value implements this interface, and if does, then
// it is used to provide the fields instead of using a generic method.
type ForEachFieldser interface {
	// ForEachField iterates through all the fields, until first false is returned.
	//
	// WARNING! This callback is supposed to be used only to read
	// the field content, but not to hold the pointer to the field itself.
	// For the purposes of performance optimizations it is assumed
	// that the pointer is not escapable. If it is required to store
	// the pointer to the field after `callback` returned then copy the field.
	// Otherwise Go's memory guarantees are broken.
	ForEachField(callback func(f *Field) bool) bool
}

// AbstractFields is an abstraction of any types of collections of fields.
type AbstractFields interface {
	ForEachFieldser
	Len() int
}

// Slice is a collection of collections of fields.
type Slice[T AbstractFields] []T

// Len implements AbstractFields.
func (s Slice[T]) Len() int {
	count := 0
	s.ForEachField(func(f *Field) bool {
		count++
		return true
	})
	return count
}

// ForEachField implements AbstractFields.
func (s Slice[T]) ForEachField(callback func(f *Field) bool) bool {
	for _, items := range s {
		if !items.ForEachField(callback) {
			return false
		}
	}
	return true
}

// Slicer implements AbstractFields providing a subset of fields (reported through
// method ForEachField of the initial collection).
type Slicer[T AbstractFields] struct {
	All      T
	StartIdx uint
	EndIdx   uint
}

var _ AbstractFields = (*Slicer[AbstractFields])(nil)

// Len implements AbstractFields.
func (s *Slicer[T]) Len() int {
	count := 0
	s.ForEachField(func(f *Field) bool {
		count++
		return true
	})
	return count
}

// ForEachField implements AbstractFields.
func (s *Slicer[T]) ForEachField(callback func(f *Field) bool) bool {
	idx := uint(0)
	return s.All.ForEachField(func(f *Field) bool {
		if idx >= s.EndIdx {
			return false
		}
		if idx >= s.StartIdx {
			if !callback(f) {
				return false
			}
		}
		idx++
		return true
	})
}

// NewSlicer provides a subset of fields (reported through method ForEachField of the initial collection).
func NewSlicer[T AbstractFields](all T, startIdx, endIdx uint) *Slicer[T] {
	return &Slicer[T]{
		All:      all,
		StartIdx: startIdx,
		EndIdx:   endIdx,
	}
}

// Gather copies all the fields into a slice and returns it.
func Gather[T AbstractFields](fields T) Fields {
	result := make(Fields, 0, fields.Len())
	fields.ForEachField(func(field *Field) bool {
		result = append(result, *field)
		return true
	})
	return result
}
