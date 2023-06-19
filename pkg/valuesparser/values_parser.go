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

package valuesparser

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/facebookincubator/go-belt/pkg/field"
)

type anyValue struct {
	Value any
}

// ForEachField implements abstract.Field, but with one difference,
// if callback is nil then the function returns true if (and only if)
// the value is a structured value (which supposed to be interpreted as
// a source of fields, but not as a part of an unstructured message).
func (s anyValue) ForEachField(callback func(f *field.Field) bool) bool {
	value := s.Value
	if value == nil {
		return true
	}

	switch v := value.(type) {
	case field.ForEachFieldser:
		if callback == nil {
			return true
		}
		return v.ForEachField(callback)
	case error:
		if callback == nil {
			return false
		}
		return callback(&field.Field{
			Key:        "error",
			Value:      v,
			Properties: nil,
		})
	}

	v := reflect.Indirect(reflect.ValueOf(value))
	switch v.Kind() {
	case reflect.Map:
		if callback == nil {
			return true
		}
		return ParseMapValue(v, callback)
	case reflect.Struct:
		if callback == nil {
			return true
		}
		return ParseStructValue(nil, v, callback)
	}

	return callback != nil
}

// AnySlice a handler for arbitrary values, which extracts structured fields/values
// by method ForEachField() and provides everything else as a string by method WriteUnparsed.
//
// It also implements fields.AbstractFields, so it could be directly used as a collection
// of structured fields right after just a type-casting ([]any -> AnySlice). It is
// supposed to be a zero allocation implementation in future.
//
// For example it is used to parse all the arguments of logger.Logger.Log.
type AnySlice []any

// ExtractUnparsed removes all values which does not define a structured field
// from the slice. The removed unstructured data is written to the writer.
func (p *AnySlice) ExtractUnparsed(w io.Writer) {
	count := 0
	for idx := 0; idx < len(*p); idx++ {
		value := (*p)[idx]
		isField := anyValue{Value: value}.ForEachField(nil)
		if isField {
			continue
		}

		fmt.Fprint(w, value)
		count++
		(*p)[idx] = nil
	}

	if count == 0 {
		return
	}

	dIdx := 0
	for sIdx := 0; sIdx < len(*p); sIdx++ {
		if (*p)[sIdx] == nil {
			continue
		}
		(*p)[dIdx] = (*p)[sIdx]
		dIdx++
	}
	*p = (*p)[:dIdx]
}

// ForEachField implements field.AbstractFields.
func (p AnySlice) ForEachField(callback func(f *field.Field) bool) bool {
	for idx := 0; idx < len(p); idx++ {
		value := p[idx]
		if !(anyValue{Value: value}.ForEachField(callback)) {
			return false
		}
	}
	return true
}

// ParseMapValue calls the callback for each pair in the map until first false is returned
//
// It returns false if callback returned false.
func ParseMapValue(m reflect.Value, callback func(f *field.Field) bool) bool {
	var f field.Field
	for _, keyV := range m.MapKeys() {
		valueV := m.MapIndex(keyV)
		switch key := keyV.Interface().(type) {
		case field.Key:
			f.Key = key
		default:
			f.Key = fmt.Sprint(key)
		}
		f.Value = valueV.Interface()
		if !callback(&f) {
			return false
		}
	}
	return true
}

// ParseStructValue parses a structure to a collection of fields.
//
// `fieldPath` is the prefix of the field-name.
// `_struct` is the structure to be parsed (provided as a reflect.Value).
// `callback` is the function called for each found field, until first false is returned.
//
// It returns false if callback returned false.
func ParseStructValue(fieldPath []string, _struct reflect.Value, callback func(f *field.Field) bool) bool {
	s := reflect.Indirect(_struct)

	var f field.Field
	// TODO: optimize this
	t := s.Type()

	fieldPath = append(fieldPath, "")

	fieldCount := s.NumField()
	for fieldNum := 0; fieldNum < fieldCount; fieldNum++ {
		structFieldType := t.Field(fieldNum)
		if structFieldType.PkgPath != "" {
			// unexported
			continue
		}
		logTag := structFieldType.Tag.Get("log")
		if logTag == "-" {
			continue
		}
		structField := s.Field(fieldNum)
		if structField.IsZero() {
			continue
		}
		value := reflect.Indirect(structField)

		pathComponent := structFieldType.Name
		if logTag != "" {
			pathComponent = logTag
		}
		fieldPath[len(fieldPath)-1] = pathComponent

		if value.Kind() == reflect.Struct {
			if !ParseStructValue(fieldPath, value, callback) {
				return false
			}
			continue
		}

		f.Key = strings.Join(fieldPath, ".")
		f.Value = value.Interface()
		if !callback(&f) {
			return false
		}
	}

	return true
}

// Len implements field.AbstractFields.
func (p AnySlice) Len() int {
	return len(p)
}
