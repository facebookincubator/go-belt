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
	"sync"

	"github.com/facebookincubator/go-belt/tool/experimental/metrics/types"
	"go.uber.org/atomic"
)

var (
	_ types.IntGauge = &IntGauge{}
)

type intGaugeFamily struct {
	sync.RWMutex
	Metrics map[string]*IntGauge
}

func (family *intGaugeFamily) get(fields AbstractFields) types.IntGauge {
	fieldsKey := fieldsToString(fields)

	family.RLock()
	metric := family.Metrics[fieldsKey]
	family.RUnlock()
	if metric != nil {
		return metric
	}

	family.Lock()
	defer family.Unlock()
	metric = family.Metrics[fieldsKey]
	if metric != nil {
		return metric
	}

	metric = &IntGauge{
		Family: family,
	}
	family.Metrics[fieldsKey] = metric

	return metric
}

// IntGauge is a naive implementation of IntGauge.
type IntGauge struct {
	Family *intGaugeFamily
	atomic.Int64
}

// Add implements metrics.IntGauge.
func (metric *IntGauge) Add(add int64) types.IntGauge {
	metric.Int64.Add(add)
	return metric
}

// Value implements metrics.IntGauge.
func (metric *IntGauge) Value() any {
	return metric.Int64.Load()
}

// WithResetFields implements IntGauge.
func (metric *IntGauge) WithResetFields(fields AbstractFields) types.IntGauge {
	return metric.Family.get(fields)
}
