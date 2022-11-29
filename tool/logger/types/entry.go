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
	"time"

	"github.com/facebookincubator/go-belt"
	"github.com/facebookincubator/go-belt/pkg/field"
)

// Entry a single log entry to be logged/written.
type Entry struct {
	// Timestamp defines the time moment when the entry was issued (for example method `Debugf` was called).
	Timestamp time.Time

	// Level is the logging level of the entry.
	Level Level

	// Message is an arbitrary string explaining the event, an unstructured message.
	// It is set by `*f` (e.g. `Debugf`) functions, by argument "message" in
	// LogFields/TraceFields/etc function or by unstructured values in Log/Trace/Debug/etc functions.
	Message string

	// Fields implements ways to read the structured fields to be logged.
	//
	// To avoid copying the fields into an intermediate format (which most likely
	// will be transformed into something else anyway in the Logger implementation)
	// we provide here an accessor to Fields instead of compiled Fields themselves.
	Fields field.AbstractFields

	// TraceIDs is the collection of unique IDs associated with this logging entry.
	// For example it may be useful to attach an unique ID to each network request,
	// so that it will be possible to fetch the whole log of processing for any network request.
	TraceIDs belt.TraceIDs

	// Properties defines special implementation-specific behavior related to the Entry.
	//
	// See description of EntryProperty.
	Properties EntryProperties

	// Caller is the Program Counter of the position in the code which initiated the logging
	// of the entry.
	//
	// See also OptionGetCallerFunc and DefaultGetCallerFunc.
	Caller PC
}

// EntryProperty defines special implementation-specific behavior related to a specific Entry.
//
// Any Emitter implementation, Hook or other tool may use it
// for internal or/and external needs.
//
// For example, a Backlogger may support both async and sync logging,
// and it could be possible to request sync logging through a property.
//
// Another example is: if an error monitor and a logger are hooked to each
// other then these properties could be used to avoid loops (to mark
// already processed entries).
type EntryProperty any

// EntryProperties is a collection of EntryProperty-es.
type EntryProperties []EntryProperty

// Has returns true if the collection of properties contains a specific EntryProperty.
// It should be equal by both: type and value.
func (s EntryProperties) Has(p EntryProperty) bool {
	for _, cmp := range s {
		if cmp == p {
			return true
		}
	}
	return false
}
