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
	"sync"

	"github.com/go-ng/sort"
	"github.com/go-ng/xsort"
)

// NewSearchableFields wraps Fields to add functionality to search
// for specific fields.
func NewSearchableFields(s Fields) SearchableFields {
	return SearchableFields{
		fields: s,
	}
}

// SearchableFields is a Fields which allows to search a specific field by key
// and deduplicate the collection.
type SearchableFields struct {
	fields      Fields
	lastSortLen uint
}

var _ AbstractFields = (*SearchableFields)(nil)

// ForEachField implements AbstractFields.
func (s *SearchableFields) ForEachField(callback func(f *Field) bool) bool {
	return s.fields.ForEachField(callback)
}

// Fields just returns Fields as they are currently stored.
//
// DO NOT MODIFY THE OUTPUT OF THIS FUNCTION UNLESS YOU
// KNOW WHAT ARE YOU DOING. IT IS NOT COPIED ONLY DUE
// TO PERFORMANCE REASONS.
func (s *SearchableFields) Fields() Fields {
	return s.fields
}

// Copy returns a copy of the SearchFields.
func (s *SearchableFields) Copy() SearchableFields {
	return s.WithPreallocate(0)
}

// WithPreallocate returns a copy of the SearchFields with the specified space preallocated
// (could be useful to append).
func (s *SearchableFields) WithPreallocate(preallocateLen uint) SearchableFields {
	return SearchableFields{fields: s.fields.WithPreallocate(preallocateLen), lastSortLen: s.lastSortLen}
}

// UnsafeOverwriteFields overwrites the collection of fields, but
// does not update the index. It could be used ONLY AND ONLY
// if the collection of fields was appended and if methods Get
// and Deduplicate were not used, yet.
//
// DO NOT USE THIS FUNCTION UNLESS YOU KNOW WHAT ARE YOU DOING.
// THIS FUNCTION IS AN IMPLEMENTATION LEAK AND WAS ADDED ONLY
// DUE TO STRONG _PERFORMANCE_ REQUIREMENTS OF PACKAGE "obs".
func (s *SearchableFields) UnsafeOverwriteFields(fields Fields) {
	if len(fields) <= len(s.fields) {
		panic("you definitely use this function wrong, stop using it!")
	}
	s.fields = fields
}

// Add just adds Field-s to the collection of Fields.
//
// The Fields are not sorted at this stage. They will
// be sorted only on first Get call.
func (s *SearchableFields) Add(fields ...Field) {
	s.fields = append(s.fields, fields...)
}

// Get returns a Field with the specific Key from the collection of Fields.
func (s *SearchableFields) Get(key Key) *Field {
	// see https://github.com/xaionaro-go/benchmarks/blob/master/search/README.md
	if len(s.fields)-int(s.lastSortLen) > 256 {
		s.sort()
	}

	if s.lastSortLen > 0 {
		// binary search
		idx := sort.Search(int(s.lastSortLen), func(i int) bool {
			return s.fields[i].Key >= key
		})
		if idx >= 0 && idx < int(s.lastSortLen) {
			return &s.fields[idx]
		}
	}

	// linear search
	for idx := len(s.fields) - 1; idx >= int(s.lastSortLen); idx-- {
		if s.fields[idx].Key == key {
			return &s.fields[idx]
		}
	}

	return nil
}

var (
	// EnableSortBufPool improves performance of Get and Deduplicate
	// methods of SearchFields, but consumes more memory.
	//
	// If you have a high performance application then you probably
	// need to enable this feature (it is unlikely to consume
	// a lot of RAM)
	EnableSortBufPool = false

	sortBufPool = sync.Pool{
		New: func() any {
			return &Fields{}
		},
	}
)

func (s *SearchableFields) sort() {
	unsortedCount := len(s.fields) - int(s.lastSortLen)
	if unsortedCount == 0 {
		return
	}
	if unsortedCount < 0 {
		panic("should not happen, ever")
	}

	// See: https://raw.githubusercontent.com/go-ng/xsort/main/BENCHMARKS.txt
	if EnableSortBufPool {
		buf := sortBufPool.Get().(*Fields)
		if cap(*buf) < unsortedCount {
			*buf = make(Fields, unsortedCount)
		} else {
			*buf = (*buf)[:unsortedCount]
		}
		xsort.AppendedWithBuf(s.fields, *buf)
		for idx := range *buf {
			(*buf)[idx] = Field{}
		}
		sortBufPool.Put(buf)
	} else {
		xsort.Appended(s.fields, uint(unsortedCount))
	}

	s.lastSortLen = uint(len(s.fields))
}

// Len returns the amount of Field-s in the collection.
func (s *SearchableFields) Len() int {
	return len(s.fields)
}

// Cap returns the capacity of the storage of the Fields.
func (s *SearchableFields) Cap() int {
	return cap(s.fields)
}

// Build sorts the fields and removes the duplicates earlier fields
// with the same key values as later fields (preserving the order)
//
// It is not recommended to call this function in high performance pieces of code,
// unless it was called shortly before that.
func (s *SearchableFields) Build() {
	if s.Len() < 2 {
		return
	}
	s.sort()

	outIdx := s.Len() - 2
	for idx := s.Len() - 2; idx >= 0; idx-- {
		prevField := &s.fields[outIdx+1]
		curField := &s.fields[idx]
		if curField.Key == prevField.Key {
			continue
		}
		if outIdx != idx {
			s.fields[outIdx] = *curField
		}
		outIdx--
	}
	s.fields = s.fields[outIdx+1:]
}
