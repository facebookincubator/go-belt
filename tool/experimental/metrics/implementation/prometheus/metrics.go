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
	"strconv"
	"strings"
	"sync"

	"github.com/facebookincubator/go-belt"
	"github.com/facebookincubator/go-belt/pkg/field"
	"github.com/facebookincubator/go-belt/tool/experimental/metrics"
	"github.com/facebookincubator/go-belt/tool/experimental/metrics/types"
	"github.com/go-ng/xsort"
	"github.com/prometheus/client_golang/prometheus"
)

var _ types.Metrics = (*Metrics)(nil)

// Fields is just a type-alias to field.Fields (for convenience).
type Fields = field.Fields

func mergeSortedStrings(dst []string, add ...string) []string {
	if len(add) == 0 {
		return dst
	}
	if len(dst) == 0 {
		return append([]string{}, add...)
	}

	dstLen := len(dst)

	if !sort.StringsAreSorted(dst) || !sort.StringsAreSorted(add) {
		panic(fmt.Sprintf("%v %v", sort.StringsAreSorted(dst), sort.StringsAreSorted(add)))
	}

	i, j := 0, 0
	for i < dstLen && j < len(add) {
		switch strings.Compare(dst[i], add[j]) {
		case -1:
			i++
		case 0:
			i++
			j++
		case 1:
			dst = append(dst, add[j])
			j++
		}
	}
	dst = append(dst, add[j:]...)
	sort.Strings(dst)
	return dst
}

func labelsWithPlaceholders(labels prometheus.Labels, placeholders []string) prometheus.Labels {
	if len(labels) == len(placeholders) {
		return labels
	}
	result := make(prometheus.Labels, len(placeholders))
	for k, v := range labels {
		result[k] = v
	}
	for _, s := range placeholders {
		if _, ok := result[s]; ok {
			continue
		}
		result[s] = ""
	}
	return result
}

// CounterVec is a family of Count metrics.
type CounterVec struct {
	*prometheus.CounterVec
	Key            string
	PossibleLabels []string
}

// AddPossibleLabels adds prometheus labels to the list of expected labels of a Metric of this family.
func (v *CounterVec) AddPossibleLabels(newLabels []string) {
	v.PossibleLabels = mergeSortedStrings(v.PossibleLabels, newLabels...)
}

// GetMetricWith returns a Metric with the values of labels as provided.
func (v *CounterVec) GetMetricWith(labels prometheus.Labels) (prometheus.Counter, error) {
	return v.CounterVec.GetMetricWith(labelsWithPlaceholders(labels, v.PossibleLabels))
}

// GaugeVec is a family of Gauge metrics.
type GaugeVec struct {
	*prometheus.GaugeVec
	Key            string
	PossibleLabels []string
}

// AddPossibleLabels adds prometheus labels to the list of expected labels of a Metric of this family.
func (v *GaugeVec) AddPossibleLabels(newLabels []string) {
	v.PossibleLabels = mergeSortedStrings(v.PossibleLabels, newLabels...)
}

// GetMetricWith returns a Metric with the values of labels as provided.
func (v *GaugeVec) GetMetricWith(labels prometheus.Labels) (prometheus.Gauge, error) {
	return v.GaugeVec.GetMetricWith(labelsWithPlaceholders(labels, v.PossibleLabels))
}

type persistentData struct {
	storage
	config
}

type storage struct {
	locker   sync.Mutex
	registry *prometheus.Registry
	count    map[string]*CounterVec
	gauge    map[string]*GaugeVec
	intGauge map[string]*GaugeVec

	// temporary buffers (placing here to avoid memory allocations)
	tmpLabels        prometheus.Labels
	tmpLabelNames    []string
	tmpLabelNamesBuf []string
}

// Metrics implements a wrapper of prometheus metrics to implement
// metrics.Metrics.
//
// Pretty slow and naive implementation. Could be improved by on-need basis.
// If you need a faster implementation, then try `tsmetrics`.
//
// Warning! This implementation does not remove automatically metrics, thus
// if a metric was created once, it will be kept in memory forever.
// If you need a version of metrics which automatically removes metrics
// non-used for a long time, then try `tsmetrics`.
//
// Warning! Prometheus does not support changing amount of labels for a metric,
// therefore we delete and create a metric from scratch if it is required to
// extend the set of labels. This procedure leads to a reset of the metric
// value.
// If you need a version of metrics which does not have such flaw, then
// try `tsmetrics` or `simplemetrics`.
type Metrics struct {
	*persistentData
	labels     prometheus.Labels
	labelNames []string
}

// New returns a new instance of Metrics
func New(registry *prometheus.Registry, opts ...Option) *Metrics {
	m := &Metrics{
		persistentData: &persistentData{
			storage: storage{
				registry:  registry,
				count:     map[string]*CounterVec{},
				gauge:     map[string]*GaugeVec{},
				intGauge:  map[string]*GaugeVec{},
				tmpLabels: make(prometheus.Labels),
			},
			config: options(opts).Config(),
		},
	}
	return m
}

var (
	alternativeDefaultRegistry     *prometheus.Registry
	alternativeDefaultRegistryOnce sync.Once
)

// Default returns metrics.Metrics with the default configuration.
var Default = func() metrics.Metrics {
	registry, ok := prometheus.DefaultRegisterer.(*prometheus.Registry)
	if !ok {
		alternativeDefaultRegistryOnce.Do(func() {
			alternativeDefaultRegistry = prometheus.NewRegistry()
		})
		registry = alternativeDefaultRegistry
	}
	return New(registry)
}

// List returns all the metrics. It could be used for a custom exporter.
func (m *Metrics) List() []prometheus.Collector {
	m.storage.locker.Lock()
	defer m.storage.locker.Unlock()
	result := make([]prometheus.Collector, 0, len(m.storage.count)+len(m.storage.gauge)+len(m.storage.intGauge))

	for _, count := range m.storage.count {
		result = append(result, count.CounterVec)
	}

	for _, gauge := range m.storage.gauge {
		result = append(result, gauge.GaugeVec)
	}

	for _, intGauge := range m.storage.intGauge {
		result = append(result, intGauge.GaugeVec)
	}

	return result
}

// Registerer returns prometheus.Registerer of this Metrics.
func (m *Metrics) Registerer() prometheus.Registerer {
	return m.registry
}

// Gatherer returns prometheus.Gatherer of this Metrics.
func (m *Metrics) Gatherer() prometheus.Gatherer {
	return m.registry
}

func (m *Metrics) getOrCreateCountVec(key string, possibleLabelNames []string) *CounterVec {
	counterVec := m.count[key]
	if counterVec != nil {
		return counterVec
	}

	counterVec = &CounterVec{
		CounterVec: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: key + "_count",
		}, possibleLabelNames),
		Key:            key,
		PossibleLabels: possibleLabelNames,
	}

	m.count[key] = counterVec

	if m.registry != nil {
		err := m.registry.Register(counterVec)
		if err != nil {
			panic(fmt.Sprintf("key: '%v', err: %v", key, err))
		}
	}

	return counterVec
}

func (m *Metrics) deleteCountVec(counterVec *CounterVec) {
	if m.registry != nil {
		if !unregister(m.registry, counterVec.CounterVec) {
			panic(counterVec)
		}
	}
	delete(m.count, counterVec.Key)
}

// Count implements context.Metrics (see the description in the interface).
func (m *Metrics) Count(key string) types.Count {
	return m.CountFields(key, nil)
}

// CountFields returns Count metric given key and additional field values (on top of already defined).
func (m *Metrics) CountFields(key string, addFields field.AbstractFields) types.Count {
	m.storage.locker.Lock()
	defer m.storage.locker.Unlock()
	labels, labelNames := m.labelsWithAddFields(addFields)

	counterVec := m.getOrCreateCountVec(key, labelNames)

	counter, err := counterVec.GetMetricWith(labels)
	if err != nil {
		m.deleteCountVec(counterVec)
		counterVec.AddPossibleLabels(labelNames)
		counterVec = m.getOrCreateCountVec(key, counterVec.PossibleLabels)
		counter, err = counterVec.GetMetricWith(labels)
		if err != nil {
			panic(err)
		}
	}

	return &Count{Metrics: m, Key: key, CounterVec: counterVec, Counter: counter}
}

// ForEachCount iterates over each Count metric. Stops on first `false` returned by the callback.
func (m *Metrics) ForEachCount(callback func(types.Count) bool) bool {
	m.storage.locker.Lock()
	defer m.storage.locker.Unlock()

	shouldExitNow := false
	for key, family := range m.storage.count {
		ch := make(chan prometheus.Metric)
		go func() {
			family.CounterVec.Collect(ch)
			close(ch)
		}()
		for abstractMetric := range ch {
			metric := abstractMetric.(prometheus.Counter)
			if !callback(&Count{Metrics: m, Key: key, CounterVec: family, Counter: metric}) {
				shouldExitNow = true
				break
			}
		}
		for range ch {
		}
		if shouldExitNow {
			return false
		}
	}
	return true
}

func (m *Metrics) getOrCreateGaugeVec(key string, possibleLabelNames []string) *GaugeVec {
	gaugeVec := m.gauge[key]
	if gaugeVec != nil {
		return gaugeVec
	}

	_gaugeVec := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: key + "_float",
	}, possibleLabelNames)

	gaugeVec = &GaugeVec{
		GaugeVec:       _gaugeVec,
		Key:            key,
		PossibleLabels: possibleLabelNames,
	}

	m.gauge[key] = gaugeVec

	if m.registry != nil {
		err := m.registry.Register(gaugeVec)
		if err != nil {
			panic(fmt.Sprintf("key: '%v', err: %v", key, err))
		}
	}

	return gaugeVec
}

func (m *Metrics) deleteGaugeVec(gaugeVec *GaugeVec) {
	if m.registry != nil {
		if !unregister(m.registry, gaugeVec.GaugeVec) {
			panic(gaugeVec)
		}
	}
	delete(m.gauge, gaugeVec.Key)
}

// expects m.storage.locker be Locked
func (m *Metrics) labelsWithAddFields(addFields field.AbstractFields) (prometheus.Labels, []string) {
	if addFields == nil {
		return m.labels, m.labelNames
	}

	for k := range m.storage.tmpLabels {
		delete(m.storage.tmpLabels, k)
	}
	for k, v := range m.labels {
		m.storage.tmpLabels[k] = v
	}
	if cap(m.storage.tmpLabelNames) < len(m.labels) {
		m.storage.tmpLabelNames = make([]string, len(m.labels)+addFields.Len())
	}
	m.storage.tmpLabelNames = m.tmpLabelNames[:len(m.labels)]
	copy(m.storage.tmpLabelNames, m.labelNames)
	addFields.ForEachField(func(f *field.Field) bool {
		if !f.Properties.Has(types.FieldPropInclude) {
			return true
		}
		v := FieldValueToString(f.Value)
		if _, isSet := m.storage.tmpLabels[f.Key]; !isSet {
			m.storage.tmpLabelNames = append(m.storage.tmpLabelNames, f.Key)
		}
		m.storage.tmpLabels[f.Key] = v
		return true
	})
	unsortedCount := len(m.storage.tmpLabelNames) - len(m.labels)
	if cap(m.storage.tmpLabelNamesBuf) < unsortedCount {
		m.storage.tmpLabelNamesBuf = make([]string, unsortedCount)
	}
	m.storage.tmpLabelNamesBuf = m.storage.tmpLabelNamesBuf[:unsortedCount]
	xsort.AppendedWithBuf(sort.StringSlice(m.storage.tmpLabelNames), m.storage.tmpLabelNamesBuf)
	return m.storage.tmpLabels, m.storage.tmpLabelNames
}

// Gauge implements context.Metrics (see the description in the interface).
func (m *Metrics) Gauge(key string) types.Gauge {
	return m.GaugeFields(key, nil)
}

// GaugeFields returns Gauge metric given key and additional field values (on top of already defined).
func (m *Metrics) GaugeFields(key string, addFields field.AbstractFields) types.Gauge {
	m.storage.locker.Lock()
	defer m.storage.locker.Unlock()
	labels, labelNames := m.labelsWithAddFields(addFields)

	gaugeVec := m.getOrCreateGaugeVec(key, labelNames)

	gauge, err := gaugeVec.GetMetricWith(labels)
	if err != nil {
		m.deleteGaugeVec(gaugeVec)
		gaugeVec.AddPossibleLabels(labelNames)
		gaugeVec = m.getOrCreateGaugeVec(key, gaugeVec.PossibleLabels)
		gauge, err = gaugeVec.GetMetricWith(labels)
		if err != nil {
			panic(err)
		}
	}

	return &Gauge{Metrics: m, Key: key, GaugeVec: gaugeVec, Gauge: gauge}
}

// ForEachGauge iterates over each Gauge metric. Stops on first `false` returned by the callback.
func (m *Metrics) ForEachGauge(callback func(types.Gauge) bool) bool {
	m.storage.locker.Lock()
	defer m.storage.locker.Unlock()

	shouldExitNow := false
	for key, family := range m.storage.gauge {
		ch := make(chan prometheus.Metric)
		go func() {
			family.GaugeVec.Collect(ch)
			close(ch)
		}()
		for abstractMetric := range ch {
			metric := abstractMetric.(prometheus.Gauge)
			if !callback(&Gauge{Metrics: m, Key: key, GaugeVec: family, Gauge: metric}) {
				shouldExitNow = true
				break
			}
		}
		for range ch {
		}
		if shouldExitNow {
			return false
		}
	}
	return true
}

func (m *Metrics) getOrCreateIntGaugeVec(key string, possibleLabelNames []string) *GaugeVec {
	gaugeVec := m.intGauge[key]
	if gaugeVec != nil {
		return gaugeVec
	}

	_gaugeVec := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: key + "_int",
	}, possibleLabelNames)

	gaugeVec = &GaugeVec{
		GaugeVec:       _gaugeVec,
		Key:            key,
		PossibleLabels: possibleLabelNames,
	}

	m.intGauge[key] = gaugeVec

	if m.registry != nil {
		err := m.registry.Register(gaugeVec)
		if err != nil {
			panic(fmt.Sprintf("key: '%v', err: %v", key, err))
		}
	}

	return gaugeVec
}

func (m *Metrics) deleteIntGaugeVec(intGaugeVec *GaugeVec) {
	if m.registry != nil {
		if !unregister(m.registry, intGaugeVec) {
			panic(intGaugeVec)
		}
	}
	delete(m.intGauge, intGaugeVec.Key)
}

// IntGauge implements context.Metrics (see the description in the interface).
func (m *Metrics) IntGauge(key string) types.IntGauge {
	return m.IntGaugeFields(key, nil)
}

// IntGaugeFields returns IntGauge metric given key and additional field values (on top of already defined).
func (m *Metrics) IntGaugeFields(key string, addFields field.AbstractFields) types.IntGauge {
	m.storage.locker.Lock()
	defer m.storage.locker.Unlock()
	labels, labelNames := m.labelsWithAddFields(addFields)

	gaugeVec := m.getOrCreateIntGaugeVec(key, labelNames)

	gauge, err := gaugeVec.GetMetricWith(labels)
	if err != nil {
		m.deleteIntGaugeVec(gaugeVec)
		gaugeVec.AddPossibleLabels(labelNames)
		gaugeVec = m.getOrCreateIntGaugeVec(key, gaugeVec.PossibleLabels)
		gauge, err = gaugeVec.GetMetricWith(labels)
		if err != nil {
			panic(err)
		}
	}

	return &IntGauge{Metrics: m, Key: key, GaugeVec: gaugeVec, Gauge: gauge}
}

// ForEachIntGauge iterates over each IntGauge metric. Stops on first `false` returned by the callback.
func (m *Metrics) ForEachIntGauge(callback func(types.IntGauge) bool) bool {
	m.storage.locker.Lock()
	defer m.storage.locker.Unlock()

	shouldExitNow := false
	for key, family := range m.storage.intGauge {
		ch := make(chan prometheus.Metric)
		go func() {
			family.GaugeVec.Collect(ch)
			close(ch)
		}()
		for abstractMetric := range ch {
			metric := abstractMetric.(prometheus.Gauge)
			if !callback(&IntGauge{Metrics: m, Key: key, GaugeVec: family, Gauge: metric}) {
				shouldExitNow = true
				break
			}
		}
		for range ch {
		}
		if shouldExitNow {
			return false
		}
	}
	return true
}

// WithContextFields implements metrics.Metrics.
func (m *Metrics) WithContextFields(allFields *field.FieldsChain, newFieldsCount int) belt.Tool {
	if m.config.DisableLabels {
		return m
	}

	result := &Metrics{
		persistentData: m.persistentData,
	}
	var (
		labels prometheus.Labels
		names  []string
	)
	if newFieldsCount == -1 {
		result.labels = nil
		result.labelNames = nil
		if allFields == nil {
			return result
		}
		labels, names = result.labelsWithAddFields(allFields)
	} else {
		labels, names = m.labelsWithAddFields(field.NewSlicer(allFields, 0, uint(newFieldsCount)))
	}
	result.labels = make(prometheus.Labels, len(labels))
	for k, v := range labels {
		result.labels[k] = v
	}
	result.labelNames = make([]string, 0, len(names))
	result.labelNames = append(result.labelNames, names...)
	return result
}

// WithTraceIDs implements metrics.Metrics.
func (m *Metrics) WithTraceIDs(traceIDs belt.TraceIDs, newTraceIDsCount int) belt.Tool {
	// Should be ignored per metrics.Metrics interface description, so returning
	// as is:
	return m
}

// Flush implements metrics.Metrics (or more specifically belt.Tool).
func (*Metrics) Flush() {}

func fieldsToLabels(fields field.AbstractFields) prometheus.Labels {
	labels := make(prometheus.Labels, fields.Len())
	fields.ForEachField(func(f *field.Field) bool {
		labels[f.Key] = FieldValueToString(f.Value)
		return true
	})
	return labels
}

const prebakeMax = 65536

var prebackedString [prebakeMax * 2]string

func init() {
	for i := -prebakeMax; i < prebakeMax; i++ {
		prebackedString[i+prebakeMax] = strconv.FormatInt(int64(i), 10)
	}
}

func getPrebakedString(v int32) string {
	if v >= prebakeMax || -v <= -prebakeMax {
		return ""
	}
	return prebackedString[v+prebakeMax]
}

// FieldValueToString converts any value to a string, which could be used
// as label value in prometheus.
func FieldValueToString(vI interface{}) string {
	switch v := vI.(type) {
	case int:
		r := getPrebakedString(int32(v))
		if len(r) != 0 {
			return r
		}
		return strconv.FormatInt(int64(v), 10)
	case uint64:
		r := getPrebakedString(int32(v))
		if len(r) != 0 {
			return r
		}
		return strconv.FormatUint(v, 10)
	case int64:
		r := getPrebakedString(int32(v))
		if len(r) != 0 {
			return r
		}
		return strconv.FormatInt(v, 10)
	case string:
		return strings.Replace(v, ",", "_", -1)
	case bool:
		switch v {
		case true:
			return "true"
		case false:
			return "false"
		}
	case []byte:
		return string(v)
	case nil:
		return "null"
	case interface{ String() string }:
		return strings.Replace(v.String(), ",", "_", -1)
	}

	return "<unknown_type>"
}
