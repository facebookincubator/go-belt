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

	"github.com/prometheus/client_golang/prometheus"
	"github.com/xaionaro-go/unsafetools"
)

// prometheus does not support full unregister. Therefore we add code
// to do so.
//
// See also: https://github.com/prometheus/client_golang/issues/203
func unregister(registerer prometheus.Registerer, c prometheus.Collector) bool {
	if !registerer.Unregister(c) {
		return false
	}

	descChan := make(chan *prometheus.Desc, 10)
	go func() {
		c.Describe(descChan)
		close(descChan)
	}()

	locker := unsafetools.FieldByName(registerer, "mtx").(*sync.RWMutex)

	dimHashesByName := *unsafetools.FieldByName(registerer, "dimHashesByName").(*map[string]uint64)

	locker.Lock()
	defer locker.Unlock()

	for desc := range descChan {
		delete(dimHashesByName, *unsafetools.FieldByName(desc, "fqName").(*string))
	}

	return true
}
