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

package tsmetrics

import (
	"github.com/facebookincubator/go-belt/tool/experimental/metrics/types"
	tsmetrics "github.com/xaionaro-go/metrics"
)

var _ types.Metrics = &Metrics{}

// Metrics implements a wrapper of github.com/xaionaro-go/metrics to implement
// metrics.Metrics.
type Metrics struct {
	Registry *tsmetrics.Registry
	Tags     *tsmetrics.FastTags
}

// New returns a new instance of Metrics
func New() *Metrics {
	m := &Metrics{
		Registry: tsmetrics.New(),
		Tags:     tsmetrics.NewFastTags().(*tsmetrics.FastTags),
	}
	m.Registry.SetDefaultGCEnabled(true)
	m.Registry.SetDefaultIsRan(true)
	m.Registry.SetSender(nil)
	return m
}

// Count implements metrics.Metrics.
func (m *Metrics) Count(key string) types.Count {
	return &Count{Metrics: m, MetricCount: m.Registry.Count(key, m.Tags)}
}

// CountFields implements metrics.Metrics.
func (m *Metrics) CountFields(key string, addFields AbstractFields) types.Count {
	tags := m.addFields(addFields)
	if tags == nil {
		return m.Count(key)
	}
	defer tags.Release()
	return &Count{Metrics: m, MetricCount: m.Registry.Count(key, tags)}
}

// ForEachCount iterates through all Count metrics. Stops at first `false` returned by the callback function.
func (m *Metrics) ForEachCount(callback func(types.Count) bool) bool {
	list := m.Registry.List()
	defer list.Release()
	for _, abstractMetric := range *list {
		metric := abstractMetric.(*tsmetrics.MetricCount)
		if !callback(&Count{Metrics: m, MetricCount: metric}) {
			return false
		}
	}
	return true
}

// Gauge implements metrics.Metrics.
func (m *Metrics) Gauge(key string) types.Gauge {
	return &Gauge{Metrics: m, MetricGaugeFloat64: m.Registry.GaugeFloat64(key, m.Tags)}
}

// GaugeFields implements metrics.Metrics.
func (m *Metrics) GaugeFields(key string, addFields AbstractFields) types.Gauge {
	tags := m.addFields(addFields)
	if tags == nil {
		return m.Gauge(key)
	}
	defer tags.Release()
	return &Gauge{Metrics: m, MetricGaugeFloat64: m.Registry.GaugeFloat64(key, tags)}
}

// ForEachGauge iterates through all Gauge metrics. Stops at first `false` returned by the callback function.
func (m *Metrics) ForEachGauge(callback func(types.Gauge) bool) bool {
	list := m.Registry.List()
	defer list.Release()
	for _, abstractMetric := range *list {
		metric := abstractMetric.(*tsmetrics.MetricGaugeFloat64)
		if !callback(&Gauge{Metrics: m, MetricGaugeFloat64: metric}) {
			return false
		}
	}
	return true
}

// IntGauge implements metrics.Metrics.
func (m *Metrics) IntGauge(key string) types.IntGauge {
	return &IntGauge{Metrics: m, MetricGaugeInt64: m.Registry.GaugeInt64(key, m.Tags)}
}

// IntGaugeFields implements metrics.Metrics.
func (m *Metrics) IntGaugeFields(key string, addFields AbstractFields) types.IntGauge {
	tags := m.addFields(addFields)
	if tags == nil {
		return m.IntGauge(key)
	}
	defer tags.Release()
	return &IntGauge{Metrics: m, MetricGaugeInt64: m.Registry.GaugeInt64(key, tags)}
}

// ForEachIntGauge iterates through all IntGauge metrics. Stops at first `false` returned by the callback function.
func (m *Metrics) ForEachIntGauge(callback func(types.IntGauge) bool) bool {
	list := m.Registry.List()
	defer list.Release()
	for _, abstractMetric := range *list {
		metric := abstractMetric.(*tsmetrics.MetricGaugeInt64)
		if !callback(&IntGauge{Metrics: m, MetricGaugeInt64: metric}) {
			return false
		}
	}
	return true
}

func setTagsFromFields(tags *tsmetrics.FastTags, fields AbstractFields) {
	fields.ForEachField(func(f *Field) bool {
		tags.Set(f.Key, f.Value)
		return true
	})
}

func (m *Metrics) addFields(addFields AbstractFields) *tsmetrics.FastTags {
	if addFields == nil {
		return nil
	}
	addLen := addFields.Len()
	if addLen == 0 {
		return nil
	}
	newTags := tsmetrics.NewFastTags().(*tsmetrics.FastTags)
	newTags.Slice = make([]*tsmetrics.FastTag, len(m.Tags.Slice), len(m.Tags.Slice)+addLen)
	copy(newTags.Slice, m.Tags.Slice)
	setTagsFromFields(newTags, addFields)
	return newTags
}

// WithContextFields implements metrics.Metrics.
func (m Metrics) WithContextFields(allFields *FieldsChain, newFieldsCount int) Tool {
	newTags := tsmetrics.NewFastTags().(*tsmetrics.FastTags)
	newTags.Slice = make([]*tsmetrics.FastTag, len(m.Tags.Slice))
	copy(newTags.Slice, m.Tags.Slice)
	count := 0
	allFields.ForEachField(func(f *Field) bool {
		count++
		if count > newFieldsCount {
			return false
		}
		if !f.Properties.Has(types.FieldPropInclude) {
			return true
		}
		newTags.Set(f.Key, f.Value)
		return true
	})
	return &Metrics{
		Registry: m.Registry,
		Tags:     newTags,
	}
}

// WithTraceIDs implements metrics.Metrics.
func (m *Metrics) WithTraceIDs(traceIDs TraceIDs, newTraceIDsCount int) Tool {
	// Should be ignored per metrics.Metrics interface description, so returning
	// as is:
	return m
}

// Flush implements metrics.Metrics (or more specifically belt.Tool).
func (*Metrics) Flush() {}
