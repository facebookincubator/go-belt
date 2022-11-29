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
	"sync"

	"github.com/facebookincubator/go-belt/pkg/field"
	"github.com/facebookincubator/go-belt/tool/experimental/metrics/types"
	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
)

var (
	_ types.IntGauge = &IntGauge{}
)

// IntGauge is an implementation of metrics.IntGauge.
type IntGauge struct {
	sync.Mutex
	*Metrics
	Key string
	*GaugeVec
	io_prometheus_client.Metric
	prometheus.Gauge
}

// Add implementations metrics.IntGauge.
func (metric *IntGauge) Add(delta int64) types.IntGauge {
	metric.Gauge.Add(float64(delta))
	metric.Lock()
	defer metric.Unlock()
	err := metric.Write(&metric.Metric)
	if err != nil {
		panic(err)
	}
	return metric
}

// Value implementations metrics.Metric.
func (metric *IntGauge) Value() any {
	metric.Lock()
	defer metric.Unlock()
	return int64(*metric.Metric.Gauge.Value)
}

// WithResetFields implements metrics.IntGauge
func (metric *IntGauge) WithResetFields(fields field.AbstractFields) types.IntGauge {
	result, err := metric.GaugeVec.GetMetricWith(fieldsToLabels(fields))
	if err == nil {
		return &IntGauge{Metrics: metric.Metrics, Key: metric.Key, GaugeVec: metric.GaugeVec, Gauge: result}
	}

	return metric.Metrics.WithContextFields(nil, -1).(types.Metrics).IntGaugeFields(metric.Key, fields)
}
