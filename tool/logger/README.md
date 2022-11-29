# Example
```go
package logger_test

import (
	"bytes"
	"context"
	"log"

	"github.com/facebookincubator/go-belt"
	"github.com/facebookincubator/go-belt/beltctx"
	"github.com/facebookincubator/go-belt/tool/logger"
	"github.com/facebookincubator/go-belt/tool/logger/implementation/glog"
	xlogrus "github.com/facebookincubator/go-belt/tool/logger/implementation/logrus"
	"github.com/facebookincubator/go-belt/tool/logger/implementation/stdlib"
	"github.com/facebookincubator/go-belt/tool/logger/implementation/zap"
	"github.com/sirupsen/logrus"
)

func Example() {
	// easy to use:
	l := xlogrus.Default()
	someFunction(1, l)

	// implementation agnostic:
	l = zap.Default()
	someFunction(2, l)

	// various implementations:
	l = glog.New()
	someFunction(3, l)

	// one may still reuse all the features of the backend logger:
	logrusInstance := logrus.New()
	logrusInstance.Formatter = &logrus.JSONFormatter{}
	l = xlogrus.New(logrusInstance)
	someFunction(4, l)

	// just another example:
	var buf bytes.Buffer
	stdLogInstance := log.New(&buf, "", log.Llongfile)
	l = stdlib.New(stdLogInstance, logger.LevelDebug)
	someFunction(5, l)

	// use go-belt to manage the logger
	obs := belt.New()
	obs = logger.BeltWithLogger(obs, l)
	someBeltyFunction(6, obs)

	// use context to manage the logger
	ctx := context.Background()
	ctx = logger.CtxWithLogger(ctx, l)
	someContextyFunction(ctx, 7)

	// use a singletony Logger:
	logger.Default = func() logger.Logger {
		return l
	}
	oneMoreFunction(8)
}

func someFunction(arg int, l logger.Logger) {
	// experience close to logrus/zap:
	l = l.WithField("arg", arg)
	anotherFunction(l)
}

func anotherFunction(l logger.Logger) {
	l.Debugf("hello world, %T!", l)
}

func someBeltyFunction(arg int, obs *belt.Belt) {
	obs = obs.WithField("arg", arg)
	anotherBeltyFunction(obs)
}

func anotherBeltyFunction(obs *belt.Belt) {
	logger.FromBelt(obs).Debugf("hello world!")
}

func someContextyFunction(ctx context.Context, arg int) {
	ctx = beltctx.WithField(ctx, "arg", arg)
	anotherContextyFunction(ctx)
}

func anotherContextyFunction(ctx context.Context) {
	logger.FromCtx(ctx).Debugf("hello world!")
}

func oneMoreFunction(arg int) {
	logger.Default().WithField("arg", arg).Debugf("hello world!")
}
```

# Interface
```go
type Logger interface {
	// Logf logs an unstructured message. Though, of course, all
	// contextual structured fields will also be logged.
	//
	// This method exists mostly for convenience, for people who
	// has not got used to proper structured logging, yet.
	// See `LogFields` and `Log`. If one have variables they want to
	// log, it is better for scalable observability to log them
	// as structured values, instead of injecting them into a
	// non-structured string.
	Logf(level Level, format string, args ...any)

	// LogFields logs structured fields with a explanation message.
	//
	// Anything that implements field.AbstractFields might be used
	// as a collection of fields to be logged.
	//
	// Examples:
	//
	// 	l.LogFields(logger.LevelDebug, "new_request", field.Fields{{Key: "user_id", Value: userID}, {Key: "group_id", Value: groupID}})
	// 	l.LogFields(logger.LevelInfo, "affected entries", field.Field{Key: "mysql_affected", Value: affectedRows})
	// 	l.LogFields(logger.LevelError, "unable to fetch user info", request) // where `request` implements field.AbstractFields
	//
	// Sometimes it is inconvenient to manually describe each field,
	// and for such cases see method `Log`.
	LogFields(level Level, message string, fields field.AbstractFields)

	// Log extracts structured fields from provided values, joins
	// the rest into an unstructured message and logs the result.
	//
	// This function provides convenience (relatively to LogFields)
	// at cost of a bit of performance.
	//
	// There are few ways to extract structured fields, which are
	// applied for each value from `values` (in descending priority order):
	// 1. If a `value` is an `*Entry` then the Entry is used (with its fields)
	// 2. If a `value` implements field.AbstractFields then ForEachField method
	//    is used (so it is become similar to LogFields).
	// 3. If a `value` is a structure (or a pointer to a structure) then
	//    fields of the structure are interpreted as structured fields
	//    to be logged (see explanation below).
	// 4. If a `value` is a map then fields a constructed out of this map.
	//
	// Everything that does not fit into any of the rules above is just
	// joined into an nonstructured message (and works the same way
	// as `message` in LogFields).
	//
	// How structures are parsed:
	// Structures are parsed recursively. Each field name of the path in a tree
	// of structures is added to the resulting field name (for example int "struct{A struct{B int}}"
	// the field name will be `A.B`).
	// To enforce another name use tag `log` (for example "struct{A int `log:"anotherName"`}"),
	// to prevent a field from logging use tag `log:"-"`.
	//
	// Examples:
	//
	// 	user, err := getUser()
	// 	if err != nil {
	// 		l.Log(logger.LevelError, err)
	// 		return err
	// 	}
	// 	l.Log(logger.LevelDebug, "current user", user) // fields of structure "user" will be logged
	// 	l.Log(logger.LevelDebug, map[string]any{"user_id": user.ID, "group_id", user.GroupID})
	// 	l.Log(logger.LevelDebug, field.Fields{{Key: "user_id", Value: user.ID}, {Key: "group_id", Value: user.GroupID}})
	// 	l.Log(logger.LevelDebug, "current user ID is ", user.ID, " and group ID is ", user.GroupID) // will result into message "current user ID is 1234 and group ID is 5678".
	Log(level Level, values ...any)

	// Emitter returns the Emitter (see the description of interface "Emitter").
	Emitter() Emitter

	// Level returns the current logging level (see description of "Level").
	Level() Level

	// TraceFields is just a shorthand for LogFields(logger.LevelTrace, ...)
	TraceFields(message string, fields field.AbstractFields)

	// DebugFields is just a shorthand for LogFields(logger.LevelDebug, ...)
	DebugFields(message string, fields field.AbstractFields)

	// InfoFields is just a shorthand for LogFields(logger.LevelInfo, ...)
	InfoFields(message string, fields field.AbstractFields)

	// WarnFields is just a shorthand for LogFields(logger.LevelWarn, ...)
	WarnFields(message string, fields field.AbstractFields)

	// ErrorFields is just a shorthand for LogFields(logger.LevelError, ...)
	ErrorFields(message string, fields field.AbstractFields)

	// PanicFields is just a shorthand for LogFields(logger.LevelPanic, ...)
	//
	// Be aware: Panic level also triggers a `panic`.
	PanicFields(message string, fields field.AbstractFields)

	// FatalFields is just a shorthand for LogFields(logger.LevelFatal, ...)
	//
	// Be aware: Panic level also triggers an `os.Exit`.
	FatalFields(message string, fields field.AbstractFields)

	// Trace is just a shorthand for Log(logger.LevelTrace, ...)
	Trace(values ...any)

	// Debug is just a shorthand for Log(logger.LevelDebug, ...)
	Debug(values ...any)

	// Info is just a shorthand for Log(logger.LevelInfo, ...)
	Info(values ...any)

	// Warn is just a shorthand for Log(logger.LevelWarn, ...)
	Warn(values ...any)

	// Error is just a shorthand for Log(logger.LevelError, ...)
	Error(values ...any)

	// Panic is just a shorthand for Log(logger.LevelPanic, ...)
	//
	// Be aware: Panic level also triggers a `panic`.
	Panic(values ...any)

	// Fatal is just a shorthand for Log(logger.LevelFatal, ...)
	//
	// Be aware: Fatal level also triggers an `os.Exit`.
	Fatal(values ...any)

	// Tracef is just a shorthand for Logf(logger.LevelTrace, ...)
	Tracef(format string, args ...any)

	// Debugf is just a shorthand for Logf(logger.LevelDebug, ...)
	Debugf(format string, args ...any)

	// Infof is just a shorthand for Logf(logger.LevelInfo, ...)
	Infof(format string, args ...any)

	// Warnf is just a shorthand for Logf(logger.LevelWarn, ...)
	Warnf(format string, args ...any)

	// Errorf is just a shorthand for Logf(logger.LevelError, ...)
	Errorf(format string, args ...any)

	// Panicf is just a shorthand for Logf(logger.LevelPanic, ...)
	//
	// Be aware: Panic level also triggers a `panic`.
	Panicf(format string, args ...any)

	// Fatalf is just a shorthand for Logf(logger.LevelFatal, ...)
	//
	// Be aware: Fatal level also triggers an `os.Exit`.
	Fatalf(format string, args ...any)

	// WithLevel returns a logger with logger level set to the given argument.
	//
	// See also the description of type "Level".
	WithLevel(Level) Logger

	// WithPreHooks returns a Logger which includes/appends pre-hooks from the arguments.
	//
	// See also description of "PreHook".
	//
	// Special case: to reset hooks use `WithHooks()` (without any arguments).
	WithPreHooks(...PreHook) Logger

	// WithHooks returns a Logger which includes/appends hooks from the arguments.
	//
	// See also description of "Hook".
	//
	// Special case: to reset hooks use `WithHooks()` (without any arguments).
	WithHooks(...Hook) Logger

	// WithField returns the logger with the added field (used for structured logging).
	WithField(key string, value any, props ...field.Property) Logger

	// WithMoreFields returns the logger with the added fields (used for structured logging)
	WithMoreFields(fields field.AbstractFields) Logger

	// WithMessagePrefix adds a string to all messages logged through the derived logger.
	WithMessagePrefix(prefix string) Logger

	// WithContextFields sets new context-defined fields. Supposed to be called
	// only by an Belt.
	//
	// allFields contains all fields as a chain of additions in a reverse-chronological order,
	// while newFieldsCount tells about how much of the fields are new (since last
	// call of WithContextFields). Thus if one will call
	// field.Slice(allFields, 0, newFieldsCount) they will get only the new fields.
	// At the same time some Tool-s may prefer just to re-set all the fields instead of adding
	// only new fields (due to performance reasons) and they may just use `allFields`.
	WithContextFields(allFields *field.FieldsChain, newFieldsCount int) Tool

	// WithTraceIDs sets new context-defined TraceIDs. Supposed to be called
	// only by an Belt.
	//
	// traceIDs and newTraceIDsCount has similar properties as allFields and newFieldsCount
	// in the WithContextFields method.
	WithTraceIDs(traceIDs TraceIDs, newTraceIDsCount int) Tool

	// Flush forces to flush all buffers.
	Flush()
}
```

# Implementations

These implementations are provided out of the box:

* [`zap`](https://pkg.go.dev/github.com/facebookincubator/go-belt/tool/logger/implementation/zap) -- is based on Uber's [`zap`](https://github.com/uber-go/zap).
* [`logrus`](https://pkg.go.dev/github.com/facebookincubator/go-belt/tool/logger/implementation/logrus) -- is based on [`github.com/sirupsen/logrus`](https://github.com/sirupsen/logrus).
* [`glog`](https://pkg.go.dev/github.com/facebookincubator/go-belt/tool/logger/implementation/glog) -- is based on Google's [`glog`](github.com/golang/glog).
* [`stdlib`](https://pkg.go.dev/github.com/facebookincubator/go-belt/tool/logger/implementation/glog) -- is based on standard Go's [`log`](https://pkg.go.dev/log) package.

# Custom implementation

Depending on how many features your logger is ready to provide it should implement one of:
* `interface{ Printf(format string, args ...any)` or just be of type `func(string, ...any)`
* [`Emitter`](https://github.com/facebookincubator/go-belt/blob/a80187cd561e4c30237aff5fccd46f06981d41e2/tool/logger/types/logger.go#L225)
* [`CompactLogger`](https://github.com/facebookincubator/go-belt/blob/a80187cd561e4c30237aff5fccd46f06981d41e2/tool/logger/adapter/compact_logger.go#L17)
* [`Logger`](https://github.com/facebookincubator/go-belt/blob/a80187cd561e4c30237aff5fccd46f06981d41e2/tool/logger/types/logger.go#L25)


And then you may call `adapter.LoggerFromAny` and it will convert your logger to `Logger` by adding everything what is missing in a naive generic way.

So the easiest implementation of a logger might be for example:
```go
import (
	"github.com/facebookincubator/go-belt/tool/logger"
	"github.com/facebookincubator/go-belt/tool/logger/adapter"
)

func getLogger() logger.Logger {
	return adapter.LoggerFromAny(func(format string, args ...any) {
		fmt.Printf("[log]"+format+"\n", args...)
	})
}
```

If a more complicated logger is required then take a look at existing [implementations](https://github.com/facebookincubator/go-belt/tree/main/tool/logger/implementation).

