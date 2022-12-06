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

package types

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/facebookincubator/go-belt/pkg/field"
	"github.com/facebookincubator/go-belt/pkg/valuesparser"
)

// Breadcrumb contains auxiliary information about something that happened before
// an Event happened, which supposed to help to investigate the Event.
type Breadcrumb struct {
	TS         time.Time
	Path       []string
	Categories []string
	Data       any
}

var _ field.ForEachFieldser = (*Breadcrumb)(nil)

// ForEachField implements field.ForEachFieldser
func (bc *Breadcrumb) ForEachField(callback func(f *field.Field) bool) bool {
	callback(&field.Field{
		Key:        "breadcrumb_ts_" + strings.Join(bc.Path, "."),
		Value:      bc.TS,
		Properties: []field.Property{fieldPropBreadcrumb},
	})
	v := reflect.Indirect(reflect.ValueOf(bc.Data))
	switch v.Kind() {
	case reflect.Struct:
		return valuesparser.ParseStructValue(bc.Path, v, callback)
	default:
		return callback(&field.Field{
			Key:        "breadcrumb_" + strings.Join(bc.Path, "."),
			Value:      fmt.Sprint(bc.Data),
			Properties: []field.Property{fieldPropBreadcrumb},
		})
	}
}
