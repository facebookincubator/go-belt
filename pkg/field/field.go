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

// Key is an identifier of the Field.
type Key = string

// Value is a value stored within a Field.
type Value = any

// Field represents a single structured value.
type Field struct {
	// Key is the identifier of the field. They are supposed to be unique,
	// but this is not a mandatory requirement. Each tool implementation
	// may treat duplicated keys differently.
	Key Key

	// Value is the value.
	Value

	// Properties defines Hook-specific options (how to interpret the specific field).
	Properties Properties
}

var _ AbstractFields = (*Field)(nil)

// ForEachField implements AbstractFields.
func (f *Field) ForEachField(callback func(f *Field) bool) bool {
	return callback(f)
}

// Len implements AbstractFields.
func (f *Field) Len() int {
	return 1
}

// Fields is a slice of Field-s
type Fields []Field

var _ AbstractFields = (Fields)(nil)

// Copy returns a copy of the slice.
func (s Fields) Copy() Fields {
	return s.WithPreallocate(0)
}

// WithPreallocate returns a copy of the slice with available additional capacity of size preallocateLen.
func (s Fields) WithPreallocate(preallocateLen uint) Fields {
	cpy := make(Fields, len(s), len(s)+int(preallocateLen))
	copy(cpy, s)
	return cpy
}

// Less implements sort.Interface
func (s Fields) Less(i, j int) bool {
	return s[i].Key < s[j].Key
}

// ForEachField implements AbstractFields
func (s Fields) ForEachField(callback func(f *Field) bool) bool {
	for idx := range s {
		if !callback(&s[idx]) {
			return false
		}
	}
	return true
}

// Len implements sort.Interface
func (s Fields) Len() int {
	return len(s)
}
