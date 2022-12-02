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
	"github.com/DataDog/gostackparse"
	"github.com/facebookincubator/go-belt/tool/experimental/tracer"
	"github.com/facebookincubator/go-belt/tool/logger"
	"github.com/google/uuid"
)

// EventID is an unique ID (usually UUIDv4) assigned to an issued event.
type EventID string

// Event is the full collection of data collected on a specific event to be reported (like a panic or an error).
type Event struct {
	// Entry is a general-purpose structured information.
	logger.Entry

	// ID is an unique ID of the Event.
	ID EventID

	// ExternalIDs is a list of IDs of the Event in external systems.
	ExternalIDs []any

	// Exception is the information about what went wrong.
	Exception

	// Spans is the tracer spans opened in the context of the Event.
	Spans tracer.Spans

	// CurrentGoroutineID is the ID of the goroutine where the Event was observed.
	CurrentGoroutineID int

	// Goroutines is the information about all the goroutines at the moment when the Event was observed.
	Goroutines []Goroutine
}

// GetID is a safe accessor of an event ID, it returns an empty ID if Event is nil.
func (ev *Event) GetID() EventID {
	if ev == nil {
		return EventID("")
	}
	return ev.ID
}

// Goroutine is a full collection of data collected on a specific goroutine.
type Goroutine = gostackparse.Goroutine

type entryPropertyErrorEvent int

const (
	// EntryPropertyErrorMonitoringEventEntry is a logger.EntryProperty implementation which marks that
	// the entry is an error monitoring event.
	EntryPropertyErrorMonitoringEventEntry = entryPropertyErrorEvent(iota + 1)

	// EntryPropertyErrorEvent is a logger.EntryProperty implementation which marks that
	// the event is an error event.
	EntryPropertyErrorEvent

	// EntryPropertyPanicEvent is a logger.EntryProperty implementation which marks that
	// the event is a panic event.
	EntryPropertyPanicEvent
)

// RandomEventID returns a new random EventID.
func RandomEventID() EventID {
	return EventID(uuid.New().String())
}
