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

package zap

import (
	"bytes"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/facebookincubator/go-belt/pkg/field"
	"github.com/facebookincubator/go-belt/tool/logger/types"
	"go.uber.org/zap"
)

type buffer struct {
	bytes.Buffer
}

func (buf *buffer) Close() error {
	return nil
}

func (buf *buffer) Sync() error {
	return nil
}

func TestLogger(t *testing.T) {
	var buf buffer
	err := zap.RegisterSink("buf", func(*url.URL) (zap.Sink, error) {
		return &buf, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	zapCfg := zap.NewDevelopmentConfig()
	zapCfg.Encoding = "json"
	zapCfg.OutputPaths = []string{"buf:"}
	zapLogger, err := zapCfg.Build()
	if err != nil {
		t.Fatal(err)
	}

	timeNow = func() time.Time {
		return time.Date(2022, 2, 24, 0, 0, 0, 0, time.UTC)
	}
	l := New(zapLogger, types.OptionGetCallerFunc(nil))
	l.ErrorFields("test", &field.Field{Key: "UserID", Value: 123})
	requireString(t, `{"L":"ERROR","T":"2022-02-24T00:00:00.000Z","M":"test","UserID":123}`, strings.Trim(buf.String(), "\n"))
	buf.Reset()

	l.Error("test", struct {
		UserID int `log:"user_id"`
	}{UserID: 123})
	requireString(t, `{"L":"ERROR","T":"2022-02-24T00:00:00.000Z","M":"test","user_id":123}`, strings.Trim(buf.String(), "\n"))
	buf.Reset()
}

func requireString(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Fatalf("expected string: '%s', actual: '%s'", expected, actual)
	}
}
