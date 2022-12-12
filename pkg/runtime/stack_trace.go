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

package runtime

import (
	"fmt"
	"runtime"
	"strings"
)

type FramesIterator interface {
	// Next works the same as standard (*runtime.Frames).Next().
	Next() (runtime.Frame, bool)
}

type StackTrace interface {
	// ProgramCounters copies a slice of program counters to the argument.
	// The program counters which begins with the final caller and ends with the initial/root caller.
	ProgramCounters(PCs) int

	// Frames is an analog of runtime.CallerFrames.
	Frames() FramesIterator

	// String implements fmt.Stringer.
	String() string

	// Len returns the amount of frames in the stack trace
	Len() int
}

// PCs is a goroutine execution frames stack trace.
type PCs []PC

var _ StackTrace = (PCs)(nil)

// Len implements interface StackTrace.
func (s PCs) Len() int {
	return len(s)
}

// ProgramCounters implements interface StackTrace.
func (s PCs) ProgramCounters(out PCs) int {
	return copy(out, s)
}

// CallersFrames is an equivalent of standard runtime.CallerFrames.
func (s PCs) Frames() FramesIterator {
	pcs := make([]uintptr, len(s))
	for idx, pc := range s {
		pcs[idx] = uintptr(pc)
	}
	return runtime.CallersFrames(pcs)
}

// String implements fmt.Stringer.
func (s PCs) String() string {
	var result strings.Builder
	frames := s.Frames()
	if frames == nil {
		return "<invalid stack trace>"
	}
	frameDepth := 1
	for {
		frame, ok := frames.Next()
		result.WriteString(fmt.Sprintf("%d. %s:%d: %s\n", frameDepth, frame.File, frame.Line, frame.Function))
		frameDepth++
		if !ok {
			break
		}
	}
	return result.String()
}

// CallerStackTrace returns the StackTrace of the current Caller (in current goroutine).
func CallerStackTrace(callerPCFilter PCFilter) PCs {
	if callerPCFilter == nil {
		callerPCFilter = DefaultCallerPCFilter
	}

	pcs := pcsPool.Get().(*[]uintptr)
	defer pcsPool.Put(pcs)

	startIdx := 0
	n := runtime.Callers(1, *pcs)
	for i := 0; i < n; i++ {
		pc := (*pcs)[i]
		if callerPCFilter(pc) {
			startIdx = i
			break
		}
	}

	result := make(PCs, n-startIdx)
	for idx, pc := range (*pcs)[startIdx:n] {
		result[idx] = PC(pc)
	}
	return result
}

// Frames is a slice of runtime.Frame
type Frames []runtime.Frame

var _ StackTrace = (Frames)(nil)

// Len implements interface StackTrace.
func (s Frames) Len() int {
	return len(s)
}

// ProgramCounters implements interface StackTrace.
func (s Frames) ProgramCounters(out PCs) int {
	endIdx := len(s)
	if len(out) < endIdx {
		endIdx = len(out)
	}
	for idx, frame := range s[:endIdx] {
		out[idx] = PC(frame.PC)
	}
	return endIdx
}

// CallersFrames is an equivalent of standard runtime.CallerFrames.
func (s Frames) Frames() FramesIterator {
	if len(s) == 0 {
		return nil
	}
	return newFramesIterator(s)
}

type framesIterator struct {
	Frames      Frames
	CurPosition int
}

func newFramesIterator(s Frames) *framesIterator {
	return &framesIterator{Frames: s}
}

var _ FramesIterator = (*framesIterator)(nil)

// Next implements FramesIterator.
func (i *framesIterator) Next() (runtime.Frame, bool) {
	if len(i.Frames) <= i.CurPosition {
		return runtime.Frame{}, false
	}

	frame := i.Frames[i.CurPosition]
	i.CurPosition++
	return frame, len(i.Frames) > i.CurPosition
}

// String implements fmt.Stringer.
func (s Frames) String() string {
	var result strings.Builder
	for idx, frame := range s {
		result.WriteString(fmt.Sprintf("%d. %s:%d: %s\n", idx+1, frame.File, frame.Line, frame.Function))
	}
	return result.String()
}
