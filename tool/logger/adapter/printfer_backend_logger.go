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
	"strconv"
	"strings"
	"sync"

	"github.com/facebookincubator/go-belt/tool/logger/types"
)

// Printfer is any object with a Printf method.
// For example the standard log.(*Logger) also implements this interface.
type Printfer interface {
	// Printf is expected to work similarly to fmt.Printf, but having
	// any output destination (not limited to stdout/stderr/whatever).
	Printf(format string, args ...any)
}

// PrintferEmitter is an implementation of a Emitter given only a Printfer.
type PrintferEmitter struct {
	Printfer
}

var _ types.Emitter = (*PrintferEmitter)(nil)

// Flush implements logger.Emitter.
func (PrintferEmitter) Flush() {}

// Emit implements logger.Emitter.
func (p PrintferEmitter) Emit(entry *types.Entry) {
	p.Printfer.Printf("%s", unstructuredMessageFromEntry(entry))
}

type printfWrap struct {
	Func func(format string, args ...any)
}

// Printf implements Printfer.
func (p printfWrap) Printf(format string, args ...any) {
	p.Func(format, args...)
}

var (
	stringsBuilderPool = sync.Pool{
		New: func() any {
			return &strings.Builder{}
		},
	}
)

func unstructuredMessageFromEntry(entry *types.Entry) string {
	result := stringsBuilderPool.Get().(*strings.Builder)
	defer func() {
		if result.Cap() >= 1024 {
			return
		}
		result.Reset()
		stringsBuilderPool.Put(result)
	}()

	msg := entry.Message
	var lineStr string

	finalLen := len(msg)
	file, line := entry.Caller.FileLine()
	if line != 0 {
		lineStr = strconv.FormatUint(uint64(line), 10)
		finalLen += len("[L :] ") + len(lineStr) + len(file)
	}
	result.Grow(finalLen)

	if line != 0 {
		result.WriteByte('[')
		result.WriteByte(entry.Level.Byte())
		result.WriteByte(' ')
		result.WriteString(file)
		result.WriteByte(':')
		result.WriteString(lineStr)
		result.WriteString("] ")
	}
	result.WriteString(msg)

	return result.String()
}
