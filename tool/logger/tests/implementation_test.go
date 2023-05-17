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

package tests

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"strings"
	"testing"

	"github.com/facebookincubator/go-belt/tool/logger/implementation/logrus"
	"github.com/facebookincubator/go-belt/tool/logger/implementation/stdlib"
	"github.com/facebookincubator/go-belt/tool/logger/implementation/zap"
	"github.com/facebookincubator/go-belt/tool/logger/types"
	upstreamlogrus "github.com/sirupsen/logrus"
	upstreamzap "go.uber.org/zap"
)

type zapBuffer struct {
	bytes.Buffer
}

func (buf *zapBuffer) Close() error {
	return nil
}

func (buf *zapBuffer) Sync() error {
	return nil
}

type logger struct {
	Name   string
	Logger types.Logger
	Output *bytes.Buffer
}

func getImplementations(t *testing.T) []logger {
	var result []logger

	// stdlib
	{
		var buf bytes.Buffer
		result = append(result, logger{
			Name:   "stdlib",
			Logger: stdlib.New(log.New(&buf, "", 0), types.LevelTrace),
			Output: &buf,
		})
	}

	// zap
	{
		var buf zapBuffer
		err := upstreamzap.RegisterSink("buf", func(*url.URL) (upstreamzap.Sink, error) {
			return &buf, nil
		})
		if err != nil {
			t.Fatal(err)
		}

		zapCfg := upstreamzap.NewDevelopmentConfig()
		zapCfg.Encoding = "json"
		zapCfg.OutputPaths = []string{"buf:"}
		zapCfg.Level = upstreamzap.NewAtomicLevelAt(zap.LevelToZap(types.LevelTrace))
		zapLogger, err := zapCfg.Build()
		if err != nil {
			t.Fatal(err)
		}
		result = append(result, logger{
			Name:   "zap",
			Logger: zap.New(zapLogger),
			Output: &buf.Buffer,
		})
	}

	// logrus
	{
		var buf bytes.Buffer
		logrusLogger := upstreamlogrus.New()
		logrusLogger.Out = &buf
		logrusLogger.Level = logrus.LevelToLogrus(types.LevelTrace)
		result = append(result, logger{
			Name:   "logrus",
			Logger: logrus.New(logrusLogger),
			Output: &buf,
		})
	}

	// glog
	{
		// the upstream glog logger does not support diverting the output to a buffer
	}

	return result
}

func TestImplementations(t *testing.T) {
	for _, l := range getImplementations(t) {
		l.Logger.Errorf("unit-test")
		l.Logger.Flush()
		if !strings.Contains(l.Output.String(), "unit-test") {
			t.Fatalf("logger %s did not print an error using Errorf", l.Name)
		}
		l.Output.Reset()

		l.Logger.Error(fmt.Errorf("unit-test"))
		l.Logger.Flush()
		if !strings.Contains(l.Output.String(), "unit-test") {
			t.Fatalf("logger %s did not print an error using Error", l.Name)
		}
		l.Output.Reset()
	}
}
