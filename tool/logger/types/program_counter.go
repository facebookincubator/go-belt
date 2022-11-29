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
	"runtime"
)

// PC is just a handy wrapper around Program Counter (PC) [see
// "program counter" in https://pkg.go.dev/runtime].
type PC uintptr

// Defined returns true if the caller is defined
func (pc PC) Defined() bool {
	return pc != 0
}

// Func returns *runtime.Func describing the caller.
func (pc PC) Func() *runtime.Func {
	if !pc.Defined() {
		return nil
	}
	return runtime.FuncForPC(uintptr(pc))
}

// FileLine returns the source code file path.
//
// In standard Go implementation this is a zero-allocation function.
func (pc PC) FileLine() (string, int) {
	if !pc.Defined() {
		return "", 0
	}
	return pc.Func().FileLine(uintptr(pc))
}

// Entry returns the entry address of the function corresponding to the program counter.
func (pc PC) Entry() uintptr {
	return pc.Func().Entry()
}
