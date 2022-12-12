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
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/facebookincubator/go-belt/pkg/runtime"
)

func TestEventGetID(t *testing.T) {
	(*Event)(nil).GetID() // should not panic
}

func TestEventError(t *testing.T) {
	ev := &Event{
		ID: "1",
		Exception: Exception{
			StackTrace: runtime.Frames{
				{
					Function: "func1()",
					File:     "file1.go",
					Line:     123,
				},
				{
					Function: "func2()",
					File:     "file2.go",
					Line:     234,
				},
			},
		},
		CurrentGoroutineID: 2,
		Goroutines: []Goroutine{
			{
				ID: 2,
			},
			{
				ID: 3,
			},
		},
	}

	t.Run("panic", func(t *testing.T) {
		ev.Exception.Error = nil
		ev.Exception.IsPanic = true
		ev.Exception.PanicValue = "unit-test"
		err := ev.AsError()
		require.ErrorAs(t, err, &ErrPanic{})
		require.Equal(t, `[event ID: '1'][goroutine 2] got a panic: unit-test
stack trace:
1. file1.go:123: func1()
2. file2.go:234: func2()
`, err.Error())
	})

	t.Run("error", func(t *testing.T) {
		ev.Exception.Error = fmt.Errorf("unit-test")
		ev.Exception.IsPanic = false
		ev.Exception.PanicValue = nil
		err := ev.AsError()
		require.ErrorAs(t, err, &ErrError{})
		require.Equal(t, "[event ID: '1'] unit-test", err.Error())
	})
}
