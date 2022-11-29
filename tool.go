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

package belt

import (
	"github.com/facebookincubator/go-belt/pkg/field"
)

// Tool is an abstract observability tool. It could be a Logger, metrics, tracing or anything else.
type Tool interface {
	// Flush forces to flush all buffers.
	Flush()

	// WithContextFields sets new context-defined fields. Supposed to be called
	// only by an Belt.
	//
	// allFields contains all fields as a chain of additions in a reverse-chronological order,
	// while newFieldsCount tells about how much of the fields are new (since last
	// call of WithContextFields). Thus if one will call
	// field.Slice(allFields, 0, newFieldsCount) they will get only the new fields.
	// At the same time some Tool-s may prefer just to re-set all the fields instead of adding
	// only new fields (due to performance reasons) and they may just use `allFields`.
	WithContextFields(allFields *field.FieldsChain, newFieldsCount int) Tool

	// WithTraceIDs sets new context-defined TraceIDs. Supposed to be called
	// only by an Belt.
	//
	// traceIDs and newTraceIDsCount has similar properties as allFields and newFieldsCount
	// in the WithContextFields method.
	WithTraceIDs(traceIDs TraceIDs, newTraceIDsCount int) Tool
}

// Tools is a collection of observability Tool-s.
type Tools map[ToolID]Tool

// GetByID returns a Tool of a specified ID. Returns an untyped nil if such Tool is not set.
func (tools Tools) GetByID(toolID ToolID) Tool {
	return tools[toolID]
}
