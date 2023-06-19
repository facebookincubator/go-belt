
<table>
<tr>
<td>

[![go report](https://goreportcard.com/badge/github.com/facebookincubator/go-belt)](https://goreportcard.com/report/github.com/facebookincubator/go-belt)
[![License](https://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)

Out of the box tools (interfaces):
|Module|GoDoc|QuickStart|
|-|-|-|
|[Logger](https://github.com/facebookincubator/go-belt/blob/main/tool/logger/types/logger.go#L37-L219)|[![GoDoc](https://godoc.org/github.com/facebookincubator/go-belt/tool/logger?status.svg)](https://pkg.go.dev/github.com/facebookincubator/go-belt/tool/logger?tab=doc)|[example](https://github.com/facebookincubator/go-belt/blob/main/tool/logger/examples/doc_test.go)|
|[Metrics](https://github.com/facebookincubator/go-belt/blob/main/tool/experimental/metrics/types/metrics.go#L20-L66) (experimental)|[![GoDoc](https://godoc.org/github.com/facebookincubator/go-belt/tool/experimental/metrics?status.svg)](https://pkg.go.dev/github.com/facebookincubator/go-belt/tool/experimental/metrics?tab=doc)|[example](https://github.com/facebookincubator/go-belt/blob/main/tool/experimental/metrics/examples/doc_test.go)|
|[Tracer](https://github.com/facebookincubator/go-belt/blob/main/tool/experimental/tracer/tracer.go#L22-L69) (experimental)|[![GoDoc](https://godoc.org/github.com/facebookincubator/go-belt/tool/experimental/tracer?status.svg)](https://pkg.go.dev/github.com/facebookincubator/go-belt/tool/experimental/tracer?tab=doc)|[example](https://github.com/facebookincubator/go-belt/blob/main/tool/experimental/tracer/examples/doc_test.go)|
|[ErrorMonitor](https://github.com/facebookincubator/go-belt/blob/main/tool/experimental/errmon/types/error_monitor.go#L47-L89) (experimental)|[![GoDoc](https://godoc.org/github.com/facebookincubator/go-belt/tool/experimental/errmon?status.svg)](https://pkg.go.dev/github.com/facebookincubator/go-belt/tool/experimental/errmon?tab=doc)|[example](https://github.com/facebookincubator/go-belt/blob/main/tool/experimental/errmon/examples/doc_test.go)|
|[Belt](https://github.com/facebookincubator/go-belt/blob/main/belt.go#L21-L34)|[![GoDoc](https://godoc.org/github.com/facebookincubator/go-belt?status.svg)](https://pkg.go.dev/github.com/facebookincubator/go-belt?tab=doc)||

Out of the box implementation examples:
|Module|Implementation|GoDoc|QuickStart|
|-|-|-|-|
|Logger|logrus|[![GoDoc](https://godoc.org/github.com/facebookincubator/go-belt/tool/logger/implementation/logrus?status.svg)](https://pkg.go.dev/github.com/facebookincubator/go-belt/tool/logger/implementation/logrus?tab=doc)|`logrus.Default()`|
|Logger|zap|[![GoDoc](https://godoc.org/github.com/facebookincubator/go-belt/tool/logger/implementation/zap?status.svg)](https://pkg.go.dev/github.com/facebookincubator/go-belt/tool/logger/implementation/zap?tab=doc)|`zap.Default()`|
|Logger|glog|[![GoDoc](https://godoc.org/github.com/facebookincubator/go-belt/tool/logger/implementation/glog?status.svg)](https://pkg.go.dev/github.com/facebookincubator/go-belt/tool/logger/implementation/glog?tab=doc)|`glog.New()`|
|Metrics|prometheus|[![GoDoc](https://godoc.org/github.com/facebookincubator/go-belt/tool/experimental/metrics/implementation/prometheus?status.svg)](https://pkg.go.dev/github.com/facebookincubator/go-belt/tool/experimental/metrics/implementation/prometheus?tab=doc)|`prometheus.Default()`|
|ErrorMonitor|sentry|[![GoDoc](https://godoc.org/github.com/facebookincubator/go-belt/tool/experimental/errmon/implementation/sentry?status.svg)](https://pkg.go.dev/github.com/facebookincubator/go-belt/tool/experimental/errmon/implementation/sentry?tab=doc)|`sentry.New(sentryClient)`|
|Tracer|zipkin|[![GoDoc](https://godoc.org/github.com/facebookincubator/go-belt/tool/experimental/tracer/implementation/zipkin?status.svg)](https://pkg.go.dev/github.com/facebookincubator/go-belt/tool/experimental/tracer/implementation/zipkin?tab=doc)|`zipkin.New(zipkinTracer)`|

</td>
<td style='vertical-align:top'>

![logo](doc/logo/variant2-small.png "logo")

</td>
<table>

# Index
1. [Mission](#mission)
2. [About](#about)
3. [Overview](#overview)
4. [Why should I use it?](#why-should-I-use-it)
5. [Quick start](#quick-start)
6. [Exotic cases](#exotic-cases)

# About

This package contains implementation-agnostic interfaces for application observability (such as logging, metrics and so on) and a collection of implementations of these interfaces. And all the observability tools are merged together into an "observability tool belt".

This project implements these ideas:
* **Dependency injection.** No observability tooling should be hardcoded into an application.
* A tool interface just represents aggregated **best practices.**
* **Easy to use.** Doing application observability comfortable in any project (is it a simple hobby project or a hyperscaler commercial project).
* Various observability tools share the idea of **structured data.** Let's take advantage of this.

This crosses with the ideas of [OpenTelemetry](https://github.com/open-telemetry/opentelemetry.io/blob/main/content/en/community/mission.md), but much more focused on "easy to use", contextual structured data, injected dependencies and deeper decoupling. Also OpenTelemetry [is more focused on metrics and tracing](https://github.com/open-telemetry/opentelemetry-go/blob/1cbd4c2b7726ded5e7a9a18997cb25977137a0b1/README.md), while this package pays more attention to logging. These projects does not compete, the opposite: for example, one may implement go-belt interfaces with OpenTelemetry SDK.

# Mission

This package intended to improve the culture of application observability in applications written in Go and to standardize approaches used across the community. In other words this is an attempt to accumulate (opinionated but:) the most generic best practices of how to handle observability (and first of all: **logging,** metrics, tracing and error monitoring).

We do not want to just propose some additional solution. We do want to collect all the best trade-offs
together and be open for changes. Please do not hesitate to propose any changes (even the most drastic ones) if you believe that will address this mission. We will try to find the best trade-offs for a generic use case and continuously improve these packages. This is the whole point of the project.

*Just in case a reminder: "the best" -- does not mean "perfect", it means "the most practical". Also more drastic the change is, more reasoning it requires.*

# Overview

There are 5 main components here:

 * [Logger](https://github.com/facebookincubator/go-belt/blob/main/tool/logger/README.md)
 * [Metrics](https://github.com/facebookincubator/go-belt/blob/main/tool/experimental/metrics/README.md)
 * [Distributed Tracer](https://github.com/facebookincubator/go-belt/blob/main/tool/experimental/tracer/README.md)
 * [Error Monitor](https://github.com/facebookincubator/go-belt/blob/main/tool/experimental/errmon/README.md)
 * and optionally the "[Belt](https://pkg.go.dev/github.com/facebookincubator/go-belt)" to access all of these and more.

All of these components are generic and abstracted from specific implementation. And some implementations are provided for each of those. For example there are implementations for [Logger](https://pkg.go.dev/github.com/facebookincubator/go-belt/tool/logger) based on: [zap](https://pkg.go.dev/github.com/facebookincubator/go-belt/tool/logger/implementation/zap), [logrus](https://pkg.go.dev/github.com/facebookincubator/go-belt/tool/logger/implementation/logrus), [glog](https://pkg.go.dev/github.com/facebookincubator/go-belt/tool/logger/implementation/glog) and [standard Go's log package](https://pkg.go.dev/github.com/facebookincubator/go-belt/tool/logger/implementation/stdlib).

And if one needs only the best practices (accumulated so far) for logging, just go to [Logger](https://pkg.go.dev/github.com/facebookincubator/go-belt/tool/logger) and disregard everything else in here. Pick an existing [implementation](https://pkg.go.dev/github.com/facebookincubator/go-belt/tool/logger/implementation) or write a new one (and create a PR to push it here). All this applies to any other tool.

But if one needs a sane control over multiple tools at the same time then use the `Belt`. This also applies to tools not mentioned here.

# Why should I use it?

* Provide **observability tooling as a dependency injection,** instead of hardcoding an implementation. For example the most of IT companies have internal infra with their-specific observability infra -- this package makes code implementation-agnostic, so that any implementation could be injected using the same shared code (see example: [ConTest](https://github.com/linuxboot/contest)).
* To make code **reusable among different projects**. To do not reinvent the same wheels over and over again working in different projects with people of different opinions about observability. In this package we try to cover the most of popular ways to do the logging, hoping it will be a good enough compromise for everybody.
* To have an **application with observability in mind**. Even if the application is already implemented without having observability tooling in mind, it is easily fixable by this package. See the [Quick start](#quick-start) section. Or if one does not plan to add proper observability at a specific moment they still can already start using this package (it does not create essential coding overhead), and in any moment in the future it will be very easy to add all the desired observability.
* To be aligned with **the best practices**. It is highly encouraged to constructively question and discuss the approaches applied here. If something is not aligned with the best practices then the goal of this project is to adapt.

# Quick start

## Logger

See also more detailed info on using `Logger` in [its README.md](https://github.com/facebookincubator/go-belt/blob/main/tool/logger/README.md).

### Approach "contextual logger"

```go
import (
	"github.com/facebookincubator/go-belt/belt"
	"github.com/facebookincubator/go-belt/tool/logger/implementation/zap"
)

func main() {
	...
	ctx = logger.CtxWithLogger(ctx, zap.Default())
	...
	someFunc(ctx)
	...
}

func someFunc(ctx context.Context) {
	...
	ctx = belt.WithField(ctx, "user_id", user.ID)
	...
	anotherFunc(ctx)
	...
}

func anotherFunc(ctx context.Context) {
	...
	logger.Debug(ctx, "hello world!") // user_id will also be logged here
	...
}
```

Also contexts usually are already propagated through a lot of codebases. Thus, one may take advantage of that in an existing codebase.

### Approach "safer contextual logger"
There is a major argument against the approach above:
* `context.Context` is pretty generic entity and may be generated by anybody, thus there could be lack of guarantee of having `Logger` in the context (because it is unclear where the specific context came from). This issue is partly mitigated through default `Logger` and default `Belt`, but in some cases there could be higher guarantee requirements.

To be sure we work with something correctly setup, it is possible to use `Belt` directly (instead of using it through a context):
```go
import (
	"github.com/facebookincubator/go-belt"
	"github.com/facebookincubator/go-belt/tool/logger/implementation/zap"
)

func main() {
	...
	belt := belt.New()
	belt = logger.BeltWithLogger(belt, zap.Default())
	...
	someFunc(ctx, belt)
	...
}

func someFunc(ctx context.Context, belt *belt.Belt) {
	...
	belt = belt.WithField("user_id", user.ID)
	...
	anotherFunc(belt)
	...
}

func anotherFunc(ctx context.Context, belt *belt.Belt) {
	...
	logger.FromBelt(belt).Debug("hello world!") // user_id will also be logged here
	...
}
```
here we always will be sure we work with the entity where Logger is correctly initialized.

### Approach "just give me structured logger"

OK-OK. Instead of injecting logger into context, you may just use it directly:
```go
// import "github.com/facebookincubator/go-belt/tool/logger/implementation/zap"

logger := zap.Default()
...
fn(logger) // the function here is agnostic of specific logger implementation
```

And this logger will be based on [Uber's zap](https://github.com/uber-go/zap).

### Approach "standard global logger"
```go
// import "github.com/facebookincubator/go-belt/tool/logger/implementation/stdlib"

stdlib.Default().Debug("Hello world!")
```
Even though this approach is discouraged, it still keeps possibility to get proper structured logging and all the fancy stuff at any moment without any difficult changes in the code.

## Metrics

Let's say we used the "contextual logger" approach for logging, and now we want to add metrics.
To do so we need only to add something like this to the initialization code:
```go
// import (
// 	promadapter "github.com/facebookincubator/go-belt/tool/experimental/metrics/implementation/prometheus"
// 	"github.com/prometheus/client_golang/prometheus"
// )
promRegistry := prometheus.NewRegistry()
ctx = metrics.CtxWithMetrics(ctx, promadapter.New(promRegistry))
```

and that's it. Now you can use prometheus metrics, for example:
```go
metrics.FromCtx(ctx).Count("requests").Add(1)
```

It will also include all the structured fields (added for example with `WithField`) allowed for metrics as labels for the prometheus metric. For example:

```go
import "github.com/facebookincubator/go-belt/tool/experimental/metrics"

func someFunc(ctx context.Context) {
	...
	ctx = belt.WithField(ctx, "user_id", user.ID, metrics.FieldPropInclude)
	...
	processRequest(ctx, req)
	...
}

func processRequest(ctx context.Context, req Request) {
	defer metrics.FromCtx(ctx).GaugeInt("concurrent_request").Add(1).Add(-1) // "user_id" will be used here as a prometheus label.
	...

	logger.FromCtx(ctx).Debug("hello world!") // and also "user_id" will be logged here as well.
}
```

It is required to add `metrics.FieldPropInclude` for fields which are used for metrics, because the amount of actual metrics proportional to the multiplication of all used values in all the labels. And some structured fields may be pretty random causing to generate unlimited amount of metrics and consume all the memory.

## Error monitor

Now let's organize application errors. It is doable through just something like:
```go
	// import "github.com/facebookincubator/go-belt/tool/experimental/errmon/implementation/sentry"
	ctx = errmon.CtxWithErrorMonitor(ctx, sentry.New(sentryClient))
```
in the initialization code, and then using the error monitor where it is required. For example:
```go
func someFunc(ctx context.Context) {
	defer func(){ errmon.ObserveRecoverCtx(ctx, recover()) }()

	...
	_, err := writer.Write(b)
	errmon.ObserveErrorCtx(ctx, err)
	...
}
```

Specifically this code will send to [Sentry](https://sentry.io/) all the errors observed from `Write` and panics in `someFunc`.

Again, all the fields (for example added through `WithField`) will also be logged in a structured way as part of the event.

Other features (like `Breadcrumb`-s) are also supported. For example:
```go
ctx = belt.WithField("breadcrumb_user.fetch", &errmon.Breadcrumb{
	TS:         time.Now()
	Path:       []string{"user", "fetch"}
	Categories: []string{"user", "network"}
	Data:       fetchUserErr,
})
```

And:
```go
errmon.ObserveErrorCtx(ctx, err)
```
now will also send the breadcrumb to the Sentry (or another error monitor implementation injected).

## Distributed tracing

Again, first initializing it:
```go
	ctx = tracer.CtxWithTracer(ctx, zipkinadapter.New(zipkinClient))
```

And use it:
```go
func mysqlQuery(ctx context.Context, query string, args ...any) {
	span, ctx := tracer.StartChildSpanFromCtx(ctx, "MySQL-query")
	defer span.Finish()
	// ..do the MySQL query here..
}
```
That's it. Of course there are more features to cover the most generic needs.

## Other tooling

Other observability tooling could be easily introduced into the `Belt`. Actually any of `logger`, `metrics`, `tracer` and `errmon` could have been provided by external projects and it would have work absolutely the same. In other words one may use these examples to create theirown standardized observability tooling and it will just work. And if you believe you created a good example of an observability tool then feel free to make a Pull Request to add it to `tool`-s here.

# More examples

See [`examples`](https://github.com/facebookincubator/go-belt/tree/main/examples) directory.

# Performance

The package is pretty much high-performance-aware despite being generic. For example using of `logrus` through this package makes it multiple times **FASTER** than using it directly. It makes even `zap` somewhat faster in a non-noop cases. This happens due to another design of handling fields, which avoids a lot of computation duplication for structured fields and provide already compiled (in a faster way) structures to the backend logger. For example on some synthetic tests it makes `zap` 15 times [faster](https://github.com/facebookincubator/go-belt/blob/main/tool/logger/implementation/zap/BENCHMARKS.txt):
```
Benchmark/prod/depth205/WithField/callLog-false/bare_zap-16       	     128	    914923 ns/op	 4883118 B/op	    3286 allocs/op
Benchmark/prod/depth205/WithField/callLog-false/adapted_zap-16    	    1870	     60388 ns/op	  103320 B/op	    1845 allocs/op
Benchmark/prod/depth205/WithField/callLog-true/bare_zap-16        	     126	    901277 ns/op	 4886755 B/op	    3290 allocs/op
Benchmark/prod/depth205/WithField/callLog-true/adapted_zap-16     	     840	    141649 ns/op	  247172 B/op	    1866 allocs/op
```

# Exotic cases

## Type-assertion

It is still allowed (though discouraged) to do type-assertion of an observability tool if it is necessary. For example:

```go
logrusEntry := logger.GetEmitter(ctx).(*logrusadapter.Emitter).LogrusEntry
logrusEntry = logrusEntry.WithFields(logrus.Fields{
	"oldFashionLogrusFieldKey": "some value",
})
logrusEntry.Debugf("hey!")
```

## Writing a custom Logger

Depending on how many features your logger is ready to provide it should implement one of:
* `interface{ Printf(format string, args ...any)` or just be of type `func(string, ...any)`
* [`Emitter`](https://github.com/facebookincubator/go-belt/blob/main/tool/logger/types/logger.go#L221-L251)
* [`CompactLogger`](https://github.com/facebookincubator/go-belt/blob/main/tool/logger/adapter/compact_logger.go#L21-L131)
* [`Logger`](https://github.com/facebookincubator/go-belt/blob/main/tool/logger/types/logger.go#L37-L219)

And then you may call `adapter.LoggerFromAny` and it will convert your logger to `Logger` by adding everything what is missing in a naive generic way.
