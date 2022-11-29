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
	"strings"
)

// FieldsChain is a chain of fields.
//
// FieldsChain is focused on high-performance Clone method.
type FieldsChain struct {
	parent *FieldsChain

	// In current implementation oneFieldKey and multipleFields are never
	// used together, but still we preserve oneFieldKey and oneFieldValue due
	// to performance reasons, it allows to avoid allocating and accessing
	// a map for a single field.
	oneField       Field
	multipleFields AbstractFields
}

var _ AbstractFields = (*FieldsChain)(nil)

// NewChainFromOne returns a FieldsChain which contains only one Field.
func NewChainFromOne(key Key, value Value, props ...Property) *FieldsChain {
	return &FieldsChain{
		oneField: Field{
			Key:        key,
			Value:      value,
			Properties: props,
		},
	}
}

// ForEachField implements AbstractFields
func (fields *FieldsChain) ForEachField(callback func(f *Field) bool) bool {
	if fields == nil {
		return true
	}

	if fields.multipleFields != nil {
		if !fields.multipleFields.ForEachField(callback) {
			return false
		}
	} else {
		if !callback(&fields.oneField) {
			return false
		}
	}

	if !fields.parent.ForEachField(callback) {
		return false
	}

	return true
}

// Len returns the total amount of fields reported by ForEachField.
// T: O(n)
//
// Implements AbstractFields
func (fields *FieldsChain) Len() int {
	length := 0
	fields.ForEachField(func(field *Field) bool {
		length++
		return true
	})
	return length
}

// WithField adds the field to the chain (from the forward -- FIFO) and
// returns the pointer to the new beginning.
func (fields *FieldsChain) WithField(key Key, value Value, props ...Property) *FieldsChain {
	return &FieldsChain{
		parent: fields,
		oneField: Field{
			Key:        key,
			Value:      value,
			Properties: props,
		},
	}
}

// WithFields adds the fields to the chain (from the forward -- FIFO) and
// returns the pointer to the new beginning.
func (fields *FieldsChain) WithFields(add AbstractFields) *FieldsChain {
	if add == nil {
		return fields
	}
	return &FieldsChain{
		parent:         fields,
		multipleFields: add,
	}
}

// WithMap adds the fields constructed from the map to the chain (from the forward -- FIFO) and
// returns the pointer to the new beginning.
func (fields *FieldsChain) WithMap(m map[Key]interface{}, props ...Property) *FieldsChain {
	if m == nil {
		return fields
	}

	if len(m) < 4 {
		result := fields
		for k, v := range m {
			result = result.WithField(k, v, props...)
		}
		return result
	}

	add := make(Fields, 0, len(m))
	for k, v := range m {
		add = append(add, Field{
			Key:        k,
			Value:      v,
			Properties: props,
		})
	}
	return fields.WithFields(add)
}

// GoString implements fmt.GoStringer
func (fields *FieldsChain) GoString() string {
	var strs []string
	fields.ForEachField(func(f *Field) bool {
		strs = append(strs, fmt.Sprintf("%#+v", *f))
		return true
	})
	return "Fields{" + strings.Join(strs, "; ") + "}"
}
