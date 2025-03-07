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

package dummy

import (
	"context"

	"github.com/facebookincubator/go-belt"
	"github.com/facebookincubator/go-belt/pkg/field"
	metrics "github.com/facebookincubator/go-belt/tool/experimental/metrics/types"
)

// Metrics is a dummy implementation of metrics.Metrics. It returns
// safe handlers which won't panic, but they doesn't do anything.
type Metrics struct{}

var _ metrics.Metrics = (*Metrics)(nil)

// New returns a new instance of Metrics.
func New() *Metrics {
	return (*Metrics)(nil)
}

// WithContextFields implements metrics.Metrics.
func (*Metrics) WithContextFields(allFields *field.FieldsChain, newFieldsCount int) belt.Tool {
	return (*Metrics)(nil)
}

// WithTraceIDs implements metrics.Metrics.
func (*Metrics) WithTraceIDs(traceIDs belt.TraceIDs, newTraceIDsCount int) belt.Tool {
	return (*Metrics)(nil)
}

// Gauge implements metrics.Metrics.
func (*Metrics) Gauge(key string) metrics.Gauge {
	return (*gauge)(nil)
}

// GaugeFields implements metrics.Metrics.
func (*Metrics) GaugeFields(key string, additionalFields field.AbstractFields) metrics.Gauge {
	return (*gauge)(nil)
}

// IntGauge implements metrics.Metrics.
func (*Metrics) IntGauge(key string) metrics.IntGauge {
	return (*intGauge)(nil)
}

// IntGaugeFields implements metrics.Metrics.
func (*Metrics) IntGaugeFields(key string, additionalFields field.AbstractFields) metrics.IntGauge {
	return (*intGauge)(nil)
}

// Count implements metrics.Metrics.
func (*Metrics) Count(key string) metrics.Count {
	return (*count)(nil)
}

// CountFields implements metrics.Metrics.
func (*Metrics) CountFields(key string, additionalFields field.AbstractFields) metrics.Count {
	return (*count)(nil)
}

// Flush implements metrics.Metrics.
func (*Metrics) Flush(context.Context) {}

type gauge struct{}

var _ metrics.Gauge = (*gauge)(nil)

// Value implements metrics.Metric.
func (*gauge) Value() any {
	return float64(0)
}

// Add implements metrics.Gauge.
func (*gauge) Add(v float64) metrics.Gauge {
	return (*gauge)(nil)
}

// WithResetFields implements metrics.Gauge.
func (*gauge) WithResetFields(field.AbstractFields) metrics.Gauge {
	return (*gauge)(nil)
}

type intGauge struct{}

var _ metrics.IntGauge = (*intGauge)(nil)

// Value implements metrics.Metric.
func (*intGauge) Value() any {
	return int64(0)
}

// Add implements metrics.IntGauge.
func (*intGauge) Add(v int64) metrics.IntGauge {
	return (*intGauge)(nil)
}

// WithResetFields implements metrics.IntGauge.
func (*intGauge) WithResetFields(field.AbstractFields) metrics.IntGauge {
	return (*intGauge)(nil)
}

type count struct{}

var _ metrics.Count = (*count)(nil)

// Value implements metrics.Metric.
func (*count) Value() any {
	return uint64(0)
}

// Add implements metrics.Count.
func (*count) Add(v uint64) metrics.Count {
	return (*count)(nil)
}

// Add implements metrics.WithResetFields.
func (*count) WithResetFields(field.AbstractFields) metrics.Count {
	return (*count)(nil)
}
