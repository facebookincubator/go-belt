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
	"github.com/facebookincubator/go-belt"
	"github.com/facebookincubator/go-belt/pkg/field"
)

// Metrics is a generic interface of a metrics handler.
//
// It implements belt.Tools, but it ignores TraceIDs reported by the Belt.
type Metrics interface {
	belt.Tool

	// Gauge returns the float64 gauge metric with key "key".
	//
	// An use case example: temperature.
	Gauge(key string) Gauge

	// GaugeFields is the same as Gauge but allows also to add
	// fields.
	//
	// In terms of Prometheus the "fields" are called "labels".
	GaugeFields(key string, additionalFields field.AbstractFields) Gauge

	// IntGauge returns the int64 gauge metric with key "key".
	//
	// An use case example: amount of concurrent requests at this moment.
	IntGauge(key string) IntGauge

	// IntGaugeFields is the same as IntGauge but allows also to add
	// fields.
	//
	// In terms of Prometheus the "fields" are called "labels".
	IntGaugeFields(key string, additionalFields field.AbstractFields) IntGauge

	// Count returns the counter metric with key "key". Count may never
	// go down. It is an monotonically increasing integer, and should
	// be interpreted that way on the backend services. For example
	// a restart of an application (which resets the metric to zero) should
	// not decrease the value on the backend.
	//
	// An use case example: total amount of requests ever received.
	Count(key string) Count

	// CountFields is the same as Count but allows also to add
	// fields.
	//
	// In terms of Prometheus the "fields" are called "labels".
	CountFields(key string, additionalFields field.AbstractFields) Count

	// TBD: extend it with percentile-oriented metrics
	// TBD: extend it with timing-oriented metrics
	// TBD: ForEach functions
}

// Metric is an abstract metric.
type Metric interface {
	// Value returns the current value of the metric
	//
	// Is not supposed to be used for anything else but metrics exporting or/and debugging/testing.
	Value() any
}

// Gauge is a float64 gauge metric.
//
// See also https://prometheus.io/docs/concepts/metric_types/
type Gauge interface {
	Metric

	// Add adds value "v" to the metric and returns the result.
	Add(v float64) Gauge

	// WithResetFields returns Gauge with fields replaces to the given ones.
	//
	// In terms of Prometheus the "fields" are called "labels".
	WithResetFields(field.AbstractFields) Gauge
}

// IntGauge is a int64 gauge metric.
//
// See also https://prometheus.io/docs/concepts/metric_types/
type IntGauge interface {
	Metric

	// Add adds value "v" to the metric and returns the result.
	Add(v int64) IntGauge

	// WithResetFields returns IntGauge with fields replaces to the given ones.
	//
	// In terms of Prometheus the "fields" are called "labels".
	WithResetFields(field.AbstractFields) IntGauge
}

// Count is a counter metric.
//
// See also https://prometheus.io/docs/concepts/metric_types/
type Count interface {
	Metric

	// Add adds value "v" to the metric and returns the result.
	Add(v uint64) Count

	// WithResetFields returns Count with fields replaces to the given ones.
	//
	// In terms of Prometheus the "fields" are called "labels".
	WithResetFields(field.AbstractFields) Count
}
