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
	"net/http"

	"github.com/facebookincubator/go-belt/pkg/field"
)

// HTTPRequest contains information about an HTTP request.
//
// Supposed to be added as a field (for example through "WithField").
type HTTPRequest http.Request

var _ field.ForEachFieldser = (*HTTPRequest)(nil)

// ForEachField implements field.ForEachFielder
func (r *HTTPRequest) ForEachField(callback func(f *field.Field) bool) bool {
	if !(field.Fields{
		{
			Key:        "request_method",
			Value:      r.Method,
			Properties: []field.Property{fieldPropRequest},
		},
		{
			Key:        "request_content_length",
			Value:      r.ContentLength,
			Properties: []field.Property{fieldPropRequest},
		},
		{
			Key:        "request_remote_addr",
			Value:      r.RemoteAddr,
			Properties: []field.Property{fieldPropRequest},
		},
	}).ForEachField(callback) {
		return false
	}

	if r.TLS != nil {
		for idx, crt := range r.TLS.PeerCertificates {
			if !callback(&field.Field{
				Key:        fmt.Sprintf("request_tls_crt_subject_%d", idx),
				Value:      crt.Subject,
				Properties: []field.Property{fieldPropRequest},
			}) {
				return false
			}

		}
	}

	if r.URL == nil {
		return true
	}

	if !callback(&field.Field{
		Key:        "request_uri",
		Value:      r.URL.RequestURI(),
		Properties: []field.Property{fieldPropRequest},
	}) {
		return false
	}

	if r.URL.User == nil {
		return true
	}

	if !callback(&field.Field{
		Key:        "request_url_user",
		Value:      r.URL.User.Username(),
		Properties: []field.Property{fieldPropRequest},
	}) {
		return false
	}

	return true
}
