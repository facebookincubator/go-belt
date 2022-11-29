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

package simplemetrics

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/facebookincubator/go-belt/tool/experimental/metrics/types"
)

var _ types.Metrics = &Metrics{}

type storage struct {
	locker           sync.RWMutex
	intGaugeFamilies map[string]*intGaugeFamily
	gaugeFamilies    map[string]*gaugeFamily
	countFamilies    map[string]*countFamily
}

// Metrics is a naive implementation of Metrics
type Metrics struct {
	*storage
	fields *FieldsChain
}

// New returns an instance of Metrics.
func New() *Metrics {
	return &Metrics{
		storage: &storage{
			intGaugeFamilies: make(map[string]*intGaugeFamily),
			gaugeFamilies:    make(map[string]*gaugeFamily),
			countFamilies:    make(map[string]*countFamily),
		},
	}
}

func cleanStringForKey(in string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(in, "=", "=="), "_", "__"), ",", "_")
}

func fieldsToString(fields AbstractFields) string {
	var fieldStrings []string
	fields.ForEachField(func(f *Field) bool {
		if !f.Properties.Has(types.FieldPropInclude) {
			return true
		}
		key := cleanStringForKey(f.Key)
		value := cleanStringForKey(fmt.Sprint(f.Value))
		fieldStrings = append(fieldStrings, fmt.Sprintf("%s=%s", key, value))
		return true
	})
	sort.Slice(fieldStrings, func(i, j int) bool {
		return fieldStrings[i] < fieldStrings[j]
	})
	return strings.Join(fieldStrings, ",")
}

// Count implements metrics.Metrics.
func (metrics *Metrics) Count(key string) types.Count {
	return metrics.CountFields(key, nil)
}

// CountFields implements metrics.Metrics.
func (metrics *Metrics) CountFields(key string, additionalFields AbstractFields) types.Count {
	fields := metrics.fields.WithFields(additionalFields)
	metrics.storage.locker.RLock()
	family := metrics.countFamilies[key]
	metrics.storage.locker.RUnlock()
	if family != nil {
		return family.get(fields)
	}

	metrics.storage.locker.Lock()
	defer metrics.storage.locker.Unlock()
	family = metrics.countFamilies[key]
	if family != nil {
		return family.get(fields)
	}

	family = &countFamily{
		Metrics: make(map[string]*Count),
	}
	metrics.countFamilies[key] = family

	return family.get(fields)
}

// ForEachCount iterates through all Count metrics. Stops at first `false` returned by the callback function.
func (metrics *Metrics) ForEachCount(callback func(types.Count) bool) bool {
	metrics.storage.locker.RLock()
	defer metrics.storage.locker.RUnlock()
	for _, family := range metrics.countFamilies {
		for _, metric := range family.Metrics {
			if !callback(metric) {
				return false
			}
		}
	}
	return true
}

// Gauge implements metrics.Metrics.
func (metrics *Metrics) Gauge(key string) types.Gauge {
	return metrics.GaugeFields(key, nil)
}

// GaugeFields implements metrics.Metrics.
func (metrics *Metrics) GaugeFields(key string, additionalFields AbstractFields) types.Gauge {
	fields := metrics.fields.WithFields(additionalFields)
	metrics.storage.locker.RLock()
	family := metrics.gaugeFamilies[key]
	metrics.storage.locker.RUnlock()
	if family != nil {
		return family.get(fields)
	}

	metrics.storage.locker.Lock()
	defer metrics.storage.locker.Unlock()
	family = metrics.gaugeFamilies[key]
	if family != nil {
		return family.get(fields)
	}

	family = &gaugeFamily{
		Metrics: make(map[string]*Gauge),
	}
	metrics.gaugeFamilies[key] = family

	return family.get(fields)
}

// ForEachGauge iterates through all Gauge metrics. Stops at first `false` returned by the callback function.
func (metrics *Metrics) ForEachGauge(callback func(types.Gauge) bool) bool {
	metrics.storage.locker.RLock()
	defer metrics.storage.locker.RUnlock()
	for _, family := range metrics.gaugeFamilies {
		for _, metric := range family.Metrics {
			if !callback(metric) {
				return false
			}
		}
	}
	return true
}

// IntGauge implements metrics.Metrics.
func (metrics *Metrics) IntGauge(key string) types.IntGauge {
	return metrics.IntGaugeFields(key, nil)
}

// IntGaugeFields implements metrics.Metrics.
func (metrics *Metrics) IntGaugeFields(key string, additionalFields AbstractFields) types.IntGauge {
	fields := metrics.fields.WithFields(additionalFields)
	metrics.storage.locker.RLock()
	family := metrics.intGaugeFamilies[key]
	metrics.storage.locker.RUnlock()
	if family != nil {
		return family.get(fields)
	}

	metrics.storage.locker.Lock()
	defer metrics.storage.locker.Unlock()
	family = metrics.intGaugeFamilies[key]
	if family != nil {
		return family.get(fields)
	}

	family = &intGaugeFamily{
		Metrics: make(map[string]*IntGauge),
	}
	metrics.intGaugeFamilies[key] = family

	return family.get(fields)
}

// ForEachIntGauge iterates through all IntGauge metrics. Stops at first `false` returned by the callback function.
func (metrics *Metrics) ForEachIntGauge(callback func(types.IntGauge) bool) bool {
	metrics.storage.locker.RLock()
	defer metrics.storage.locker.RUnlock()
	for _, family := range metrics.intGaugeFamilies {
		for _, metric := range family.Metrics {
			if !callback(metric) {
				return false
			}
		}
	}
	return true
}

// WithContextFields implements metrics.Metrics.
func (metrics Metrics) WithContextFields(allFields *FieldsChain, newFieldsCount int) Tool {
	metrics.fields = allFields
	return &metrics
}

// WithTraceIDs implements metrics.Metrics.
func (metrics *Metrics) WithTraceIDs(traceIDs TraceIDs, newTraceIDsCount int) Tool {
	// Should be ignored per metrics.Metrics interface description, so returning
	// as is:
	return metrics
}

// Flush implements metrics.Metrics (or more specifically belt.Tool).
func (*Metrics) Flush() {}
