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

package errmon

import (
	"github.com/facebookincubator/go-belt/tool/experimental/errmon/types"
)

// ErrorMonitor is an observability Tool (belt.Tool) which allows
// to report about any exceptions which happen for debugging. It
// collects any useful information it can.
//
// An ErrorMonitor implementation is not supposed to be fast, but
// it supposed to provide verbose reports (sufficient enough to
// debug found problems).
type ErrorMonitor = types.ErrorMonitor

// Event is the full collection of data collected on a specific event to be reported (like a panic or an error).
type Event = types.Event

// Breadcrumb contains auxiliary information about something that happened before
// an Event happened, which supposed to help to investigate the Event.
//
// Supposed to be added as a field (for example through "WithField").
type Breadcrumb = types.Breadcrumb

// User contains information about an user.
//
// Supposed to be added as a field (for example through "WithField").
type User = types.User

// HTTPRequest contains information about an HTTP request.
//
// Supposed to be added as a field (for example through "WithField").
type HTTPRequest = types.HTTPRequest

// Package contains information about a software package.
//
// Supposed to be added as a field (for example through "WithField").
type Package = types.Package

// Tag contains information about a tag.
//
// Supposed to be added as a field (for example through "WithField").
type Tag = types.Tag
