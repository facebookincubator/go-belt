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

// Copyright (c) Facebook, Inc. and its affiliates.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package tests

import (
	"testing"

	"github.com/facebookincubator/go-belt/pkg/field"
	"github.com/facebookincubator/go-belt/tool/experimental/metrics/types"
	"github.com/stretchr/testify/require"
)

// Fields is a type alias to field.Fields (for convenience)
type Fields = field.Fields

// FieldsChain is a type alias to field.FieldsChain (for convenience)
type FieldsChain = field.FieldsChain

// FieldProperties is a type alias to field.Properties (for convenience)
type FieldProperties = field.Properties

// testMetric tests metric and returns the resulting value of the metric
func testMetric(
	t *testing.T,
	m types.Metrics,
	key string,
	fields Fields,
	resetFields Fields,
	metricType string,
	expectedValue float64,
) float64 {
	cFields := (*FieldsChain)(nil).WithFields(fields)
	switch metricType {
	case "Count":
		metric := m.WithContextFields(cFields, cFields.Len()).(types.Metrics).Count(key)
		if resetFields != nil {
			metric = metric.WithResetFields(resetFields)
		}
		require.Equal(t, uint64(expectedValue+1), metric.Add(1).Value().(uint64))
		metric = m.WithContextFields(cFields, cFields.Len()).(types.Metrics).Count(key)
		if resetFields != nil {
			metric = metric.WithResetFields(resetFields)
		}
		require.Equal(t, uint64(expectedValue+2), metric.Add(1).Value().(uint64))
		scope := m
		if resetFields != nil {
			scope = scope.WithContextFields(nil, 0).(types.Metrics)
			fields = resetFields
		}
		metric = scope.CountFields(key, fields)
		require.Equal(t, uint64(expectedValue+3), metric.Add(1).Value().(uint64))
		return float64(metric.Value().(uint64))
	case "Gauge":
		metric := m.WithContextFields(cFields, cFields.Len()).(types.Metrics).Gauge(key)
		if resetFields != nil {
			metric = metric.WithResetFields(resetFields)
		}
		require.Equal(t, expectedValue-1.5, metric.Add(-1.5).Value().(float64))
		metric = m.WithContextFields(cFields, cFields.Len()).(types.Metrics).Gauge(key)
		if resetFields != nil {
			metric = metric.WithResetFields(resetFields)
		}
		require.Equal(t, expectedValue-0.5, metric.Add(1).Value().(float64))
		scope := m
		if resetFields != nil {
			scope = scope.WithContextFields(nil, 0).(types.Metrics)
			fields = resetFields
		}
		metric = scope.GaugeFields(key, fields)
		require.Equal(t, expectedValue+0.5, metric.Add(1).Value().(float64))
		require.Equal(t, expectedValue, metric.Add(-0.5).Value().(float64))
		return metric.Value().(float64)
	case "IntGauge":
		metric := m.WithContextFields(cFields, cFields.Len()).(types.Metrics).IntGauge(key)
		if resetFields != nil {
			metric = metric.WithResetFields(resetFields)
		}
		require.Equal(t, int64(-1), metric.Add(-1).Value().(int64))
		metric = m.WithContextFields(cFields, cFields.Len()).(types.Metrics).IntGauge(key)
		if resetFields != nil {
			metric = metric.WithResetFields(resetFields)
		}
		require.Equal(t, int64(0), metric.Add(1).Value().(int64))
		scope := m
		if resetFields != nil {
			scope = scope.WithContextFields(nil, 0).(types.Metrics)
			fields = resetFields
		}
		metric = scope.IntGaugeFields(key, fields)
		require.Equal(t, int64(2), metric.Add(2).Value().(int64))
		require.Equal(t, int64(0), metric.Add(-2).Value().(int64))
		return float64(metric.Value().(int64))
	}

	panic(metricType)
}

// TestMetrics performs generic tests for an abstract types.Metrics implementation.
func TestMetrics(t *testing.T, metricsFactory func() types.Metrics) {
	m := metricsFactory()

	for _, metricType := range []string{"Count", "Gauge", "IntGauge"} {
		t.Run(metricType, func(t *testing.T) {
			t.Run("hello_world", func(t *testing.T) {
				testMetric(t, m, "hello_world", nil, nil, metricType, 0)
				testMetric(t, m, "hello_world", Fields{
					{Key: "testField", Value: "testValue", Properties: FieldProperties{types.FieldPropInclude}},
				}, nil, metricType, 0)
				testMetric(t, m, "hello_world", Fields{
					{Key: "testField", Value: "anotherValue", Properties: FieldProperties{types.FieldPropInclude}},
				}, nil, metricType, 0)
				testMetric(t, m, "hello_world", Fields{
					{Key: "anotherField", Value: "testValue", Properties: FieldProperties{types.FieldPropInclude}},
				}, nil, metricType, 0)
				testMetric(t, m, "hello_world", Fields{
					{Key: "testField", Value: "testValue", Properties: FieldProperties{types.FieldPropInclude}},
					{Key: "anotherField", Value: "anotherValue", Properties: FieldProperties{types.FieldPropInclude}},
				}, nil, metricType, 0)
			})
			t.Run("WithResetFields", func(t *testing.T) {
				key := "WithResetFields"
				tags := Fields{
					{Key: "testField", Value: "testValue", Properties: FieldProperties{types.FieldPropInclude}},
				}
				prevResult := testMetric(t, m, key, tags, nil, metricType, 0)

				wrongTags := Fields{
					{Key: "testField", Value: "anotherValue", Properties: FieldProperties{types.FieldPropInclude}},
				}
				prevResult = testMetric(t, m, key, wrongTags, tags, metricType, prevResult)

				wrongTags = Fields{
					{Key: "anotherField", Value: "anotherValue", Properties: FieldProperties{types.FieldPropInclude}},
				}
				testMetric(t, m, key, wrongTags, tags, metricType, prevResult)
			})
		})
	}
}
