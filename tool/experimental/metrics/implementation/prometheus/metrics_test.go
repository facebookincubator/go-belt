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

package prometheus

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/facebookincubator/go-belt/pkg/field"
	"github.com/facebookincubator/go-belt/tool/experimental/metrics"
	metricstester "github.com/facebookincubator/go-belt/tool/experimental/metrics/tests"
	"github.com/facebookincubator/go-belt/tool/experimental/metrics/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func contextFields(fields map[string]interface{}) *field.FieldsChain {
	var cFields *field.FieldsChain
	for k, v := range fields {
		cFields = cFields.WithField(k, v, metrics.AllowInMetrics)
	}
	return cFields
}

func TestMetrics(t *testing.T) {
	for name, registerer := range map[string]*prometheus.Registry{
		"WithRegisterer":    prometheus.NewRegistry(),
		"WithoutRegisterer": nil,
	} {
		t.Run(name, func(t *testing.T) {
			metricstester.TestMetrics(t, func() types.Metrics {
				cFields := contextFields(map[string]interface{}{"testField": "", "anotherField": ""})

				// Current implementation resets metrics if new label appears,
				// thus some unit-tests fails (and the should). Specifically
				// for prometheus we decided to tolerate this problem, therefore
				// adding hacks to prevent a wrong values: pre-register metrics
				// with all the labels beforehand.
				m := New(registerer)
				m.WithContextFields(cFields, -1).(types.Metrics).Count("WithResetFields")
				m.WithContextFields(cFields, -1).(types.Metrics).Gauge("WithResetFields")
				m.WithContextFields(cFields, -1).(types.Metrics).IntGauge("WithResetFields")

				return m
			})
		})
	}
}

func TestMetricsList(t *testing.T) {
	m := New(nil)
	c0 := m.WithContextFields(contextFields(map[string]interface{}{"testField": "", "anotherField": ""}), -1).(types.Metrics).Count("test")
	g0 := m.WithContextFields(contextFields(map[string]interface{}{"testField": "", "anotherField": ""}), -1).(types.Metrics).Gauge("test")
	i0 := m.WithContextFields(contextFields(map[string]interface{}{"testField": "", "anotherField": ""}), -1).(types.Metrics).IntGauge("test")
	c1 := m.WithContextFields(contextFields(map[string]interface{}{"testField": "a", "anotherField": ""}), -1).(types.Metrics).Count("test")
	g1 := m.WithContextFields(contextFields(map[string]interface{}{"testField": "b", "anotherField": ""}), -1).(types.Metrics).Gauge("test")
	i1 := m.WithContextFields(contextFields(map[string]interface{}{"testField": "c", "anotherField": ""}), -1).(types.Metrics).IntGauge("test")

	list := m.List()
	require.Len(t, list, 3)
	for _, c := range list {
		ch := make(chan prometheus.Metric)
		go func() {
			c.Collect(ch)
			close(ch)
		}()

		count := 0
		for range ch {
			count++
		}
		assert.Equal(t, 2, count, fmt.Sprintf("%#+v", c))
	}
	if t.Failed() {
		t.Errorf("c0: %#+v\ng0: %#+v\ni0: %#+v\nc1: %#+v\ng1: %#+v\ni1: %#+v\n",
			c0, g0, i0, c1, g1, i1)
	}
}

func TestMetricsRegistererDoubleUse(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics0 := New(registry)
	metrics1 := New(registry)

	// these test cases should panic:

	t.Run("Count", func(t *testing.T) {
		defer func() {
			require.NotNil(t, recover())
		}()

		metrics0.Count("test")
		metrics1.Count("test")
	})

	t.Run("Gauge", func(t *testing.T) {
		defer func() {
			require.NotNil(t, recover())
		}()

		metrics0.Gauge("test")
		metrics1.Gauge("test")
	})

	t.Run("IntGauge", func(t *testing.T) {
		defer func() {
			require.NotNil(t, recover())
		}()

		metrics0.IntGauge("test")
		metrics1.IntGauge("test")
	})
}

func TestMergeSortedStrings(t *testing.T) {
	slices := [][]string{
		{"a", "b", "c"},
		{"r", "a", "n", "d", "o", "m"},
		{"a", "rb", "", "it", "r", "ary"},
	}
	for idx := range slices {
		sort.Strings(slices[idx])
	}
	for _, a := range slices {
		for _, b := range slices {
			t.Run(strings.Join(a, "-")+"_"+strings.Join(b, "-"), func(t *testing.T) {
				m := map[string]struct{}{}
				for _, aItem := range a {
					m[aItem] = struct{}{}
				}
				for _, bItem := range b {
					m[bItem] = struct{}{}
				}
				expected := make([]string, 0, len(m))
				for k := range m {
					expected = append(expected, k)
				}
				sort.Strings(expected)

				require.Equal(t, expected, mergeSortedStrings(a, b...))
			})
		}
	}
}

func TestDisabledLabels(t *testing.T) {
	registry := prometheus.NewRegistry()
	var m metrics.Metrics = New(registry, OptionDisableLabels(true))
	m = m.WithContextFields(((*field.FieldsChain)(nil)).WithField("1", 2), 1).(metrics.Metrics)
	c := m.Count("someCount")
	require.Len(t, c.(*Count).labelNames, 0)
	c = c.WithResetFields(field.Fields{{Key: "3", Value: 4}})
	require.Len(t, c.(*Count).labelNames, 0)
}
