// Copyright 2023 Meta Platforms, Inc. and affiliates.
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
	"context"

	"github.com/facebookincubator/go-belt/internal"
	"github.com/facebookincubator/go-belt/pkg/field"
)

type ctxKeyType = internal.BeltCtxKeyType

var ctxKey = ctxKeyType{}

// CtxWithBelt returns a context derivative which includes the Belt as a value.
func CtxWithBelt(ctx context.Context, belt *Belt) context.Context {
	return context.WithValue(ctx, ctxKey, belt)
}

// CtxBelt returns the Belt from context values. Returns the default observer if one is not set in the context.
func CtxBelt(ctx context.Context) *Belt {
	observer := ctx.Value(ctxKey)
	if observer == nil {
		return Default()
	}
	return observer.(*Belt)
}

// WithField returns a context with a clone/derivative of the Belt which includes the passed value.
//
// The value is used by observability tooling. For example a Logger derived from the resulting
// Belt may add this value to the structured fields of each log entry.
func WithField(ctx context.Context, key string, value interface{}, props ...field.Property) context.Context {
	return context.WithValue(ctx, ctxKey, CtxBelt(ctx).WithField(key, value, props...))
}

// WithFields is the same as WithField, but adds multiple Fields at the same time.
//
// It is more performance efficient than adding fields by one.
func WithFields(ctx context.Context, fields field.AbstractFields) context.Context {
	return context.WithValue(ctx, ctxKey, CtxBelt(ctx).WithFields(fields))
}

// WithMap is just a sugar method, which provides logrus like way of adding fields.
// Effectively the same as WithFields, just the argument are in another format.
func WithMap(ctx context.Context, m map[string]interface{}, props ...field.Property) context.Context {
	return context.WithValue(ctx, ctxKey, CtxBelt(ctx).WithMap(m, props...))
}

// WithTool returns a context with an Belt clone/derivative, but the provided tool
// added to the collection of tools.
//
// Special case: to remove a specific tool, just passed an untyped nil as `tool`.
func WithTool(ctx context.Context, toolID ToolID, tool Tool) context.Context {
	return context.WithValue(ctx, ctxKey, CtxBelt(ctx).WithTool(toolID, tool))
}

// WithTraceID returns a context with an Belt clone/derivative with the passed traceIDs added to the set of TraceIDs.
func WithTraceID(ctx context.Context, traceIDs ...TraceID) context.Context {
	return context.WithValue(ctx, ctxKey, CtxBelt(ctx).WithTraceID(traceIDs...))
}

// WithArtifact returns a derivative of the context, but with the Artifact set.
func WithArtifact(ctx context.Context, artifactID ArtifactID, artifact Artifact) context.Context {
	return context.WithValue(ctx, ctxKey, CtxBelt(ctx).WithArtifact(artifactID, artifact))
}

// GetFields returns returns the set of fields set in the scope of this Belt.
//
// Do not modify the output of this function! It is for reading only.
func GetFields(ctx context.Context) field.AbstractFields {
	return CtxBelt(ctx).Fields()
}

// GetArtifacts returns the collection of Artifacts in the scope of the Belt.
//
// Do not modify the output of this function! It is for reading only.
func GetArtifacts(ctx context.Context) Artifacts {
	return CtxBelt(ctx).Artifacts()
}

// GetTraceIDs returns the current set of TraceID-s.
//
// Do not modify the output of this function! It is for reading only.
func GetTraceIDs(ctx context.Context) TraceIDs {
	return CtxBelt(ctx).TraceIDs()
}

// GetTools returns the current collection of Tools.
//
// Do not modify the output of this function! It is for reading only.
func GetTools(ctx context.Context) Tools {
	return CtxBelt(ctx).Tools()
}

// CtxFlush forces to flush all buffers of all the tools.
func Flush(ctx context.Context) {
	CtxBelt(ctx).Flush(ctx)
}
