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

package errmon

import (
	"context"

	"github.com/facebookincubator/go-belt/beltctx"
)

// FromCtx returns an ErrorMonitor given a context.
//
// Returns the Default ErrorMonitor if one is not set.
func FromCtx(ctx context.Context) ErrorMonitor {
	return FromBelt(beltctx.Belt(ctx))
}

// CtxWithErrorMonitor returns a context with an ErrorMonitor added.
func CtxWithErrorMonitor(ctx context.Context, errMon ErrorMonitor) context.Context {
	return beltctx.WithTool(ctx, ToolID, errMon)
}

// ObserveErrorCtx calls ObserveError method of an ErrorMonitor in the given context.
//
// If one at any moment in the code has an error they want to report about if it is not nil,
// they can just use this function.
//
// For example:
//
//	_, err := writer.Write(b)
//	errmon.ObserveErrorCtx(ctx, err)
func ObserveErrorCtx(ctx context.Context, err error) *Event {
	belt := beltctx.Belt(ctx)
	return FromBelt(belt).ObserveError(belt, err)
}

// ObserveRecoverCtx calls ObserveRecover method of an ErrorMonitor in the given context.
//
// If one at any moment in the code has a potential panic they want to report about,
// they can just use this function.
//
// For example one may add this to a beginning of a function:
//
//	defer func() { errmon.ObserveRecoverCtx(ctx, recover()) }()
//
// See also: https://go.dev/ref/spec#Handling_panics
func ObserveRecoverCtx(ctx context.Context, recoverResult any) *Event {
	belt := beltctx.Belt(ctx)
	return FromBelt(belt).ObserveRecover(belt, recoverResult)
}
