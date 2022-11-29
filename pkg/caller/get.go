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

package caller

import (
	"runtime"
	"strings"
	"sync"
)

var (
	pcsPool = sync.Pool{
		New: func() any {
			pcs := make([]uintptr, 1000)
			return &pcs
		},
	}
)

// PC returns the Program Counter (PC) of the caller.
//
// The caller in the call stack is defined by the frameFilter.
// If the frameFilter is nil then the DefaultPCFilter is used instead.
func PC(frameFilter PCFilter) uintptr {
	if frameFilter == nil {
		frameFilter = DefaultPCFilter
	}

	pcs := pcsPool.Get().(*[]uintptr)
	defer pcsPool.Put(pcs)

	n := runtime.Callers(1, *pcs)
	for i := 0; i < n; i++ {
		pc := (*pcs)[i]
		if frameFilter(pc) {
			return pc
		}
	}

	return 0
}

// PCFilter is a function which returns false if a Program Counter (PC) is not a caller.
//
// It is used to filter out calls inside the observability tooling and provide the real caller.
//
// The function is called sequentially from the top of the call stack until first true is met.
type PCFilter func(pc uintptr) bool

// DefaultPCFilter is an overridable function used to get the default PCFilter, used
// by function PC, if one was not provided.
var DefaultPCFilter PCFilter = func(pc uintptr) bool {
	fn := runtime.FuncForPC(pc)
	funcName := fn.Name()
	switch {
	case strings.Contains(funcName, "github.com/facebookincubator/go-belt"),
		strings.HasPrefix(funcName, "runtime"):
		return false
	}
	return true
}
