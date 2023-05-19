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
	"sync"

	"github.com/facebookincubator/go-belt/pkg/field"
)

// Belt (tool belt for observability) is the handler which orchestrates
// all the observability tooling.
//
// Things like Logger, Tracer, Metrics, ErrorMonitoring are handled together,
// and are available at any moment, given an Belt.
//
// For example one may call `logger.FromBelt(belt)` to get a `Logger`, if `belt` is an
// `Belt`. It is different from just storing this tooling as values within a context,
// for example it is:
// * More performance efficient.
// * Safer (e.g. you won't get an untyped nil because a context was accidentally changed to a new one)
// * More convenient: no need to manage each tool separately. And all the possible sugar is already in place.
// * Reusable. It is a generic standardized way to handle observability tooling.
// * Quality. Here we try to accumulate best practices, instead of having a quick solution for a single project.
type Belt struct {
	haveTools        bool
	tools            Tools
	toolsLocker      sync.Mutex
	artifacts        Artifacts
	contextFields    *field.FieldsChain
	newFieldsCount   int
	traceIDs         TraceIDs
	newTraceIDsCount int
}

// New returns a new instance of an Belt
func New() *Belt {
	return &Belt{}
}

func (belt *Belt) clone() *Belt {
	return &Belt{
		haveTools:        belt.haveTools,
		tools:            belt.tools,
		artifacts:        belt.artifacts,
		traceIDs:         belt.traceIDs,
		contextFields:    belt.contextFields,
		newFieldsCount:   belt.newFieldsCount,
		newTraceIDsCount: belt.newTraceIDsCount,
	}
}

// WithField returns a clone/derivative of the Belt which includes the passed value.
//
// The value is used by observability tooling. For example a Logger derived from the resulting
// Belt may add this value to the structured fields of each log entry.
func (belt *Belt) WithField(key string, value field.Value, props ...field.Property) *Belt {
	clone := belt.clone()
	clone.contextFields = clone.contextFields.WithField(key, value, props...)
	clone.newFieldsCount++
	return clone
}

// WithFields is the same as WithField, but adds multiple Fields at the same time.
//
// It is more performance efficient than adding fields by one.
func (belt *Belt) WithFields(fields field.AbstractFields) *Belt {
	clone := belt.clone()
	clone.contextFields = clone.contextFields.WithFields(fields)
	clone.newFieldsCount += fields.Len()
	return clone
}

// WithMap is just a sugar method, which provides logrus like way of adding fields.
// Effectively the same as WithFields, just the argument are in another format.
func (belt *Belt) WithMap(m map[string]interface{}, props ...field.Property) *Belt {
	clone := belt.clone()
	clone.contextFields = clone.contextFields.WithMap(m, props...)
	clone.newFieldsCount += len(m)
	return clone
}

// Fields returns returns the set of fields set in the scope of this Belt.
//
// Do not modify the output of this function! It is for reading only.
func (belt *Belt) Fields() field.AbstractFields {
	return belt.contextFields
}

// WithTool returns an Belt clone/derivative, but the provided tool
// added to the collection of tools.
//
// Special case: to remove a specific tool, just passed an untyped nil as `tool`.
func (belt *Belt) WithTool(toolID ToolID, tool Tool) *Belt {
	clone := belt.clone()

	tools := make(Tools, len(clone.tools)+1)
	belt.toolsLocker.Lock()
	for k, v := range clone.tools {
		tools[k] = v
	}
	belt.toolsLocker.Unlock()
	if tool == nil {
		delete(tools, toolID)
	} else {
		tools[toolID] = tool
	}
	clone.tools = tools
	clone.haveTools = len(tools) > 0

	return clone
}

// Tools returns the current collection of Tools.
//
// Do not modify the output of this function! It is for reading only.
func (belt *Belt) Tools() Tools {
	if !belt.haveTools {
		return nil
	}
	belt.toolsLocker.Lock()
	defer belt.toolsLocker.Unlock()
	belt.updateTools()
	return belt.tools
}

func (belt *Belt) updateTools() {
	if belt.newFieldsCount == 0 && belt.newTraceIDsCount == 0 {
		return
	}

	newTools := make(Tools, len(belt.tools))
	for toolID, tool := range belt.tools {
		if belt.newFieldsCount != 0 {
			tool = tool.WithContextFields(belt.contextFields, belt.newFieldsCount)
		}
		if belt.newTraceIDsCount != 0 {
			tool = tool.WithTraceIDs(belt.traceIDs, belt.newTraceIDsCount)
		}
		newTools[toolID] = tool
	}
	belt.tools = newTools

	belt.newFieldsCount = 0
	belt.newTraceIDsCount = 0
}

// WithArtifact returns a clone of the Belt, but with the Artifact set.
func (belt *Belt) WithArtifact(artifactID ArtifactID, artifact Artifact) *Belt {
	clone := belt.clone()

	artifacts := make(Artifacts, len(clone.artifacts))
	for k, v := range clone.artifacts {
		artifacts[k] = v
	}
	if artifact == nil {
		delete(artifacts, artifactID)
	} else {
		artifacts[artifactID] = artifact
	}
	clone.artifacts = artifacts

	return clone
}

// Artifacts returns the collection of Artifacts in the scope of the Belt.
//
// Do not modify the output of this function! It is for reading only.
func (belt *Belt) Artifacts() Artifacts {
	return belt.artifacts
}

// WithTraceID returns an Belt clone/derivative with the passed traceIDs added to the set of TraceIDs.
func (belt *Belt) WithTraceID(traceIDs ...TraceID) *Belt {
	clone := belt.clone()
	clone.traceIDs = append(belt.traceIDs, traceIDs...)
	clone.newTraceIDsCount += len(traceIDs)
	return clone
}

// TraceIDs returns the current set of TraceID-s.
//
// Do not modify the output of this function! It is for reading only.
func (belt *Belt) TraceIDs() TraceIDs {
	return belt.traceIDs
}

// Flush forces to flush all buffers of all the tools.
func (belt *Belt) Flush() {
	for _, tool := range belt.Tools() {
		tool.Flush()
	}
}
