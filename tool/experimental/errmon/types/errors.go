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
)

// ErrPanic is a wrapper of an Event which implements `error` in case if
// the Event conveys a panic event.
type ErrPanic struct {
	Event *Event
}

// Error implements interface `error`.
func (err ErrPanic) Error() string {
	return fmt.Sprintf(`[event ID: '%s'][goroutine %d] got a panic: %v
stack trace:
%s`,
		err.Event.ID,
		err.Event.CurrentGoroutineID,
		err.Event.PanicValue,
		err.Event.StackTrace,
	)
}

// ErrError is a wrapper of an Event which implements `error` in case if
// the Event conveys an error event.
type ErrError struct {
	Event *Event
}

// Error implements interface `error`.
func (err ErrError) Error() string {
	return fmt.Sprintf("[event ID: '%s'] %v", err.Event.ID, err.Unwrap())
}

// Unwrap implements go1.13 errors unwrapping.
func (err ErrError) Unwrap() error {
	return err.Event.Exception.Error
}
