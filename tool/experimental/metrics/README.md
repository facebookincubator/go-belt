# Disclaimer

This API is experimental and has no stability guarantees.

# Example
```go
package metrics_test

import (
	"context"

	"github.com/facebookincubator/go-belt"
	"github.com/facebookincubator/go-belt/beltctx"
	"github.com/facebookincubator/go-belt/tool/experimental/metrics"
	"github.com/facebookincubator/go-belt/tool/experimental/metrics/implementation/prometheus"
	"github.com/facebookincubator/go-belt/tool/experimental/metrics/implementation/tsmetrics"
)

func Example() {
	var m metrics.Metrics

	// easy to use:
	m = prometheus.Default()
	someFunction(1, m)

	// implementation agnostic:
	m = tsmetrics.New()
	someFunction(2, m)

	// one may still reuse all the features of the backend Metrics:
	m.(*tsmetrics.Metrics).Registry.SetDisabled(true)

	// use go-belt to manage the Metrics
	obs := belt.New()
	obs = metrics.BeltWithMetrics(obs, m)
	someBeltyFunction(3, obs)

	// use context to manage the Metrics
	ctx := context.Background()
	ctx = metrics.CtxWithMetrics(ctx, m)
	someContextyFunction(ctx, 4)

	// use a singletony Metrics:
	metrics.Default = func() metrics.Metrics {
		return m
	}
	yetOneMoreFunction(5)
}

func someFunction(arg int, m metrics.Metrics) {
	// experience close to logrus/zap:
	m = metrics.WithField(m, "arg", arg)
	anotherFunction(m)
}

func anotherFunction(m metrics.Metrics) {
	m.Count("hello").Add(1)
}

func someBeltyFunction(arg int, obs *belt.Belt) {
	obs = obs.WithField("arg", arg)
	anotherBeltyFunction(obs)
}

func anotherBeltyFunction(obs *belt.Belt) {
	metrics.FromBelt(obs).Count("hello").Add(1)
}

func someContextyFunction(ctx context.Context, arg int) {
	ctx = beltctx.WithField(ctx, "arg", arg)
	anotherContextyFunction(ctx)
}

func anotherContextyFunction(ctx context.Context) {
	metrics.FromCtx(ctx).Count("hello").Add(1)
}

func yetOneMoreFunction(arg int) {
	m := metrics.Default()
	m = metrics.WithField(m, "arg", arg)
	m.Count("hello").Add(1)
}
```

# Interface
```go
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
	// be interpreted that way on the emitter services. For example
	// a restart of an application (which resets the metric to zero) should
	// not decrease the value on the emitter.
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
```

