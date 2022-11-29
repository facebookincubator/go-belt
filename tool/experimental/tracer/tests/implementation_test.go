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

package tests

import (
	"context"
	"testing"

	"github.com/facebookincubator/go-belt/tool/experimental/tracer"
)

type DummyReporter interface {
	OnSend(func(span tracer.Span))
}

type implementationCase struct {
	Name    string
	Factory func() (tracer.Tracer, DummyReporter)
}

var implementations []implementationCase

func TestImplementations(t *testing.T) {
	for _, implCase := range implementations {
		t.Run(implCase.Name, func(t *testing.T) {
			t.Run("case:parent-child-finish-finish", func(t *testing.T) {
				tracerImpl, reporter := implCase.Factory()

				checkParentChild := func(t *testing.T, parent, child tracer.Span) {
					if tracer.IsNoopSpan(parent) {
						t.Fatalf("the parent Span is NOOP")
					}
					if tracer.IsNoopSpan(child) {
						t.Fatalf("the child Span is NOOP")
					}
					sendCount := 0
					reporter.OnSend(func(span tracer.Span) {
						sendCount++
						if span.Name() != "child" {
							t.Fatalf("unexpected name: %s", span.Name())
						}
					})
					child.Finish()
					child.Flush()
					reporter.OnSend(func(span tracer.Span) {
						sendCount++
						if span.Name() != "parent" {
							t.Fatalf("unexpected name: %s", span.Name())
						}
					})
					parent.Finish()
					parent.Flush()
					if sendCount != 2 {
						t.Fatalf("expected sendCount is 2, but got %d", sendCount)
					}
				}

				t.Run("without_ctx", func(t *testing.T) {
					parent := tracerImpl.Start("parent", nil)
					child := tracerImpl.Start("child", parent)
					checkParentChild(t, parent, child)
				})

				t.Run("with_ctx", func(t *testing.T) {
					ctx := context.Background()
					parent, ctx := tracerImpl.StartWithCtx(ctx, "parent")
					child, ctx := tracerImpl.StartChildWithCtx(ctx, "child")
					_ = ctx
					checkParentChild(t, parent, child)
				})
			})
		})
	}
}
