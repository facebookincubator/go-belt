// Copyright 2023 Meta Platforms, Inc. and affiliates.
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

package formatter

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var logLevelSymbol []byte

func init() {
	logLevelSymbol = make([]byte, len(logrus.AllLevels)+1)
	for _, level := range logrus.AllLevels {
		logLevelSymbol[level] = strings.ToUpper(level.String()[:1])[0]
	}
}

// CompactText is a logrus formatter which prints laconic lines, like
// [2001-02-03T04:05:06Z W main.go:56] my message
type CompactText struct {
	TimestampFormat string
	FieldAllowList  []string
}

// Format implements logrus.Formatter.
func (f *CompactText) Format(entry *logrus.Entry) ([]byte, error) {
	var str, header strings.Builder
	timestamp := time.RFC3339
	if f.TimestampFormat != "" {
		timestamp = f.TimestampFormat
	}
	header.WriteString(fmt.Sprintf("%s %c",
		entry.Time.Format(timestamp),
		logLevelSymbol[entry.Level],
	))
	if entry.Caller != nil {
		header.WriteString(fmt.Sprintf(" %s:%d", filepath.Base(entry.Caller.File), entry.Caller.Line))
	}
	str.WriteString(fmt.Sprintf("[%s] %s",
		header.String(),
		entry.Message,
	))

	keys := make([]string, 0, len(entry.Data))
	for key := range entry.Data {
		if f.FieldAllowList != nil {
			found := false
			for _, allowed := range f.FieldAllowList {
				if key == allowed {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		str.WriteString(fmt.Sprintf("\t%s=%v", key, entry.Data[key]))
	}

	str.WriteByte('\n')
	return []byte(str.String()), nil
}
