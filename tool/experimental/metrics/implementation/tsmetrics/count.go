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

var (
	_ types.Count = &Count{}
)

// Count is an implementation of metrics.Count.
type Count struct {
	*Metrics
	*tsmetrics.MetricCount
}

// Add implementations metrics.Count.
func (metric *Count) Add(delta uint64) types.Count {
	metric.MetricCount.Add(int64(delta))
	return metric
}

// Value implements metrics.Count.
func (metric *Count) Value() any {
	return uint64(metric.MetricCount.Get())
}

// WithResetFields implements metrics.Count.
func (metric *Count) WithResetFields(fields AbstractFields) types.Count {
	tags := tsmetrics.NewFastTags().(*tsmetrics.FastTags)
	setTagsFromFields(tags, fields)
	key := metric.MetricCount.GetName()
	return &Count{
		Metrics:     metric.Metrics,
		MetricCount: metric.Metrics.Registry.Count(key, tags),
	}
}
