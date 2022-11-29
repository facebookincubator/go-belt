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

package adapter

import (
	"bytes"
	"runtime"

	"github.com/DataDog/gostackparse"

	errmontypes "github.com/facebookincubator/go-belt/tool/experimental/errmon/types"
)

func getGoroutines() ([]errmontypes.Goroutine, int) {
	// TODO: consider pprof.Lookup("goroutine") instead of runtime.Stack

	// getting all goroutines
	stackBufferSize := 65536 * runtime.NumGoroutine()
	if stackBufferSize > 10*1024*1024 {
		stackBufferSize = 10 * 1024 * 1024
	}
	stackBuffer := make([]byte, stackBufferSize)
	n := runtime.Stack(stackBuffer, true)
	goroutines, errs := gostackparse.Parse(bytes.NewReader(stackBuffer[:n]))
	if len(errs) > 0 { //nolint:staticcheck
		// TODO: do something
	}

	// convert goroutines for the output
	goroutinesConverted := make([]errmontypes.Goroutine, 0, len(goroutines))
	for _, goroutine := range goroutines {
		goroutinesConverted = append(goroutinesConverted, *goroutine)
	}

	// getting current goroutine ID
	n = runtime.Stack(stackBuffer, false)
	currentGoroutines, errs := gostackparse.Parse(bytes.NewReader(stackBuffer[:n]))
	if len(errs) > 0 { //nolint:staticcheck
		// TODO: do something
	}
	var currentGoroutineID int
	switch len(currentGoroutines) {
	case 0:
		// TODO: do something
	case 1:
		currentGoroutineID = currentGoroutines[0].ID
	default:
		// TODO: do something
	}

	return goroutinesConverted, currentGoroutineID
}
