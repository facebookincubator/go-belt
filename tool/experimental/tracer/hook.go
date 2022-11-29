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

package tracer

// Hook is executed right before sending a Span (to some backend).
type Hook interface {
	// ProcessSpan is the main method of a Hook. It is executed before sending
	// a Span (to some backend). If it returns false then any further processing
	// and sending of the Span is cancelled.
	ProcessSpan(Span) bool

	// Flush just gracefully empties all the queues (if there are any).
	Flush()
}

// Hooks is a collection of multiple Hook-s.
type Hooks []Hook

var _ Hook = (Hooks)(nil)

// ProcessSpan implements Hook.
func (s Hooks) ProcessSpan(span Span) bool {
	for _, hook := range s {
		if !hook.ProcessSpan(span) {
			return false
		}
	}
	return true
}

// Flush implements Hook.
func (s Hooks) Flush() {
	for _, hook := range s {
		hook.Flush()
	}
}
