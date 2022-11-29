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

package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/facebookincubator/go-belt"
	"github.com/facebookincubator/go-belt/beltctx"
	"github.com/facebookincubator/go-belt/pkg/field"
	"github.com/facebookincubator/go-belt/tool/experimental/errmon"
	errmonlogger "github.com/facebookincubator/go-belt/tool/experimental/errmon/implementation/logger"
	"github.com/facebookincubator/go-belt/tool/experimental/metrics"
	prometheusadapter "github.com/facebookincubator/go-belt/tool/experimental/metrics/implementation/prometheus"
	"github.com/facebookincubator/go-belt/tool/experimental/tracer"
	"github.com/facebookincubator/go-belt/tool/experimental/tracer/implementation/zipkin"
	"github.com/facebookincubator/go-belt/tool/logger"
	"github.com/facebookincubator/go-belt/tool/logger/implementation/zap"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
)

func main() {
	ctx := context.Background()

	ctx = logger.CtxWithLogger(ctx, zap.Default())
	ctx = metrics.CtxWithMetrics(ctx, prometheusadapter.New(prometheus.NewRegistry()))
	ctx = tracer.CtxWithTracer(ctx, zipkin.Default())
	ctx = errmon.CtxWithErrorMonitor(ctx, errmonlogger.New(logger.FromCtx(ctx)))

	ctx = beltctx.WithTraceID(ctx, belt.RandomTraceID())
	ctx = beltctx.WithField(ctx, "pid", os.Getpid(), metrics.AllowInMetrics)
	ctx = beltctx.WithFields(ctx, &field.Field{Key: "field0", Value: "value0"})

	go func() {
		errmon.ObserveErrorCtx(ctx, http.ListenAndServe("localhost:6060", nil))
	}()

	doSomething(ctx)
}

func doSomething(ctx context.Context) {
	span, ctx := tracer.StartSpanFromCtx(ctx, "doSomething")
	defer span.Finish()
	// prints:
	// 2022/07/06 11:45:53 2022-07-06 11:45:53.513474862 +0100 IST m=+0.010575838:
	// {
	//   "timestamp": 1657104353512656,
	//   "duration": 798,
	//   "traceId": "0bab5ae6fef17cc9",
	//   "id": "0bab5ae6fef17cc9",
	//   "name": "dosomething",
	//   "tags": {
	//     "field0": "value0",
	//     "pid": "1662371",
	//     "trace_ids": "d8cfd7a8-2b28-4b9a-9848-5a4a481cecbf"
	//   }
	// }

	ctx = beltctx.WithField(ctx, "uid", os.Getuid())
	logger.FromCtx(ctx).Infof("yay!")
	// prints:
	// {"level":"info","ts":1657286934.0679588,"msg":"yay!","uid":1000,"field0":"value0","pid":2446918}

	logFields(ctx)
	useLogger(beltctx.WithField(ctx, "someField", "someValue"))
	tryChildSpan(ctx)
	tryPanic(ctx)
	observerErr(ctx)
	useMetrics(ctx)
}

func tryPanic(ctx context.Context) {
	defer func() { errmon.ObserveRecoverCtx(ctx, recover()) }()
	// prints:
	// {"level":"error","ts":1657286934.0684156,"caller":"everything/main.go:77","msg":"got panic with argument: some well-respected reason to panic","uid":1000,"field0":"value0","pid":2446918,"trace_id":["5e4e44ce-272f-4e0f-b66c-6a7dff511ce6"]}

	panic("some well-respected reason to panic")
}

func observerErr(ctx context.Context) {
	for _, err := range []error{fmt.Errorf("some error"), nil} {
		errmon.ObserveErrorCtx(ctx, err)
	}
	// prints:
	// {"level":"error","ts":1657286934.068724,"caller":"everything/main.go:86","msg":"some error","uid":1000,"field0":"value0","pid":2446918,"trace_id":["5e4e44ce-272f-4e0f-b66c-6a7dff511ce6"]}
}

func logFields(ctx context.Context) {
	log := logger.FromCtx(ctx)

	beltctx.Belt(ctx).Fields().ForEachField(func(f *field.Field) bool {
		log.Debugf("%#+v", *f)
		return true
	})
	// prints:
	// {"level":"debug","ts":1657286934.068012,"msg":"field.Field{Key:\"uid\", Value:1000, Properties:field.Properties(nil)}","uid":1000,"field0":"value0","pid":2446918}
	// {"level":"debug","ts":1657286934.0680254,"msg":"field.Field{Key:\"field0\", Value:\"value0\", Properties:field.Properties(nil)}","uid":1000,"field0":"value0","pid":2446918}
	// {"level":"debug","ts":1657286934.0680408,"msg":"field.Field{Key:\"pid\", Value:2446918, Properties:field.Properties{types.fieldPropertyAllowInMetricsT{}}}","uid":1000,"field0":"value0","pid":2446918}
}

func useLogger(ctx context.Context) {
	logger.FromCtx(ctx).InfoFields("Hello here as well!", &field.Field{Key: "custom_field", Value: 123})
	// prints:
	// {"level":"info","ts":1657286934.0680518,"msg":"Hello here as well!","uid":1000,"field0":"value0","pid":2446918,"someField":"someValue","custom_field":123}
}

func tryChildSpan(ctx context.Context) {
	span, ctx := tracer.StartChildSpanFromCtx(ctx, "tryChildSpan")
	defer span.Finish()
	// prints:
	// 2022/07/08 14:28:54 2022-07-08 14:28:54.068304102 +0100 IST m=+0.011247634:
	// {
	//   "timestamp": 1657286934068076,
	//   "duration": 1,
	//   "traceId": "1e832b96d001606b",
	//   "id": "0b8352f99252595c",
	//   "parentId": "1e832b96d001606b",
	//   "name": "trychildspan",
	//   "tags": {
	//     "field0": "value0",
	//     "pid": "2446918",
	//     "trace_ids": "5e4e44ce-272f-4e0f-b66c-6a7dff511ce6",
	//     "uid": "1000"
	//   }
	// }

	// doing something here
	_ = ctx
}

func useMetrics(ctx context.Context) {
	defer metrics.FromCtx(ctx).IntGauge("concurrent").Add(1).Add(-1)
	metrics.FromCtx(ctx).Count("total").Add(1)

	exportMetrics(ctx)
}

func exportMetrics(ctx context.Context) {
	// exporting metrics is not an implementation-agnostic procedure
	metrics, err := metrics.FromCtx(ctx).(*prometheusadapter.Metrics).Gatherer().Gather()
	errmon.ObserveErrorCtx(ctx, err)

	var buf bytes.Buffer
	enc := expfmt.NewEncoder(&buf, expfmt.FmtText)
	for _, metric := range metrics {
		if err := enc.Encode(metric); err != nil {
			errmon.ObserveErrorCtx(ctx, err)
			return
		}
	}

	fmt.Print(buf.String())
	// prints:
	// # HELP concurrent_int
	// # TYPE concurrent_int gauge
	// concurrent_int{pid="2446918"} 1
	// # HELP total_count
	// # TYPE total_count counter
	// total_count{pid="2446918"} 1
}
