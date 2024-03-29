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

package logger

import (
	"github.com/facebookincubator/go-belt"
	"github.com/facebookincubator/go-belt/tool/logger/types"
)

// Entry is just a type-alias for logger/types.Entry for convenience.
type Entry = types.Entry

// Hook is just a type-alias for logger/types.Hook for convenience.
type Hook = types.Hook

// Hooks is just a type-alias for logger/types.Hooks for convenience.
type Hooks = types.Hooks

// PreHook is just a type-alias for logger/types.PreHook for convenience.
type PreHook = types.PreHook

// PreHooks is just a type-alias for logger/types.PreHooks for convenience.
type PreHooks = types.PreHooks

// Logger is just a type-alias for logger/types.Logger for convenience.
type Logger = types.Logger

// Emitter is just a type-alias for logger/types.Emitter for convenience.
type Emitter = types.Emitter

// EntryProperty is just a type-alias for logger/types.EntryProperty for convenience.
type EntryProperty = types.EntryProperty

var _ belt.Tool = (Logger)(nil)
