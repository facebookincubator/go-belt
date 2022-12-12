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

package sentry

import (
	"fmt"
	"strings"

	"github.com/facebookincubator/go-belt/pkg/field"
	"github.com/facebookincubator/go-belt/pkg/runtime"
	"github.com/facebookincubator/go-belt/tool/experimental/errmon"
	errmonadapter "github.com/facebookincubator/go-belt/tool/experimental/errmon/adapter"
	errmontypes "github.com/facebookincubator/go-belt/tool/experimental/errmon/types"
	"github.com/facebookincubator/go-belt/tool/experimental/tracer"
	loggertypes "github.com/facebookincubator/go-belt/tool/logger/types"
	"github.com/getsentry/sentry-go"
)

// Emitter is a wrapper for a Sentry client to implement errmon.Emitter.
type Emitter struct {
	SentryClient *sentry.Client
}

// NewEmitter returns a new instance of Emitter.
func NewEmitter(sentryClient *sentry.Client) *Emitter {
	return &Emitter{
		SentryClient: sentryClient,
	}
}

// New wraps a Sentry client and returns a new instance, which implements errmon.ErrorMonitor.
func New(
	sentryClient *sentry.Client,
	opts ...Option,
) errmon.ErrorMonitor {
	return errmonadapter.ErrorMonitorFromEmitter(
		NewEmitter(sentryClient),
		options(opts).Config().CallerFrameFilter,
	)
}

// Flush implements errmon.Emitter
func (*Emitter) Flush() {}

// Emit implements errmon.Emitter.
func (h *Emitter) Emit(ev *errmon.Event) {
	sendEvent := EventToSentry(ev)

	eventID := h.SentryClient.CaptureEvent(sendEvent, nil, nil)
	if eventID == nil {
		return
	}

	ev.ExternalIDs = append(ev.ExternalIDs, eventID)
}

// LevelToSentry returns the closest sentry analog of a given logger.Level
func LevelToSentry(level loggertypes.Level) sentry.Level {
	switch level {
	case loggertypes.LevelTrace, loggertypes.LevelDebug:
		return sentry.LevelDebug
	case loggertypes.LevelInfo:
		return sentry.LevelInfo
	case loggertypes.LevelWarning:
		return sentry.LevelWarning
	case loggertypes.LevelError:
		return sentry.LevelError
	case loggertypes.LevelPanic, loggertypes.LevelFatal:
		return sentry.LevelFatal
	default:
		return sentry.LevelError
	}
}

// FuncNameToSentryModule converts a funcation name (see runtime.Frame) to
// a sentry module name.
func FuncNameToSentryModule(funcName string) string {
	return strings.Split(funcName, ".")[0]
}

// GoroutinesToSentry converts goroutines to the Sentry format.
func GoroutinesToSentry(goroutines []errmontypes.Goroutine, currentGoroutineID int) []sentry.Thread {
	result := make([]sentry.Thread, 0, len(goroutines))
	for _, goroutine := range goroutines {
		converted := sentry.Thread{
			ID:      fmt.Sprint(goroutine.ID),
			Current: goroutine.ID == currentGoroutineID,
			Stacktrace: &sentry.Stacktrace{
				Frames: make([]sentry.Frame, 0, len(goroutine.Stack)),
			},
		}
		if goroutine.LockedToThread {
			converted.Name = fmt.Sprintf("goroutine_lockedToThread_%08X", goroutine.ID)
		} else {
			converted.Name = fmt.Sprintf("goroutine_%08X", goroutine.ID)
		}
		for _, frame := range goroutine.Stack {
			if frame.Func == "panic" && strings.HasSuffix(frame.File, "runtime/panic.go") {
				converted.Crashed = true
			}

			converted.Stacktrace.Frames = append(converted.Stacktrace.Frames, sentry.Frame{
				Function: frame.Func,
				Filename: frame.File,
				Lineno:   frame.Line,
				Module:   FuncNameToSentryModule(frame.Func),
			})
		}

		result = append(result, converted)
	}
	return result
}

// StackTraceToSentry converts a stack trace to the Sentry format.
func StackTraceToSentry(stackTrace runtime.StackTrace) *sentry.Stacktrace {
	frames := stackTrace.Frames()
	if frames == nil {
		return nil
	}
	result := &sentry.Stacktrace{
		Frames: make([]sentry.Frame, 0, stackTrace.Len()),
	}
	for {
		frame, ok := frames.Next()
		result.Frames = append(result.Frames, sentry.Frame{
			Function: frame.Function,
			Module:   FuncNameToSentryModule(frame.Function),
			Filename: frame.File,
			Lineno:   frame.Line,
		})
		if !ok {
			break
		}
	}
	return result
}

// SpansToSentry converts tracer spans to the Sentry format.
func SpansToSentry(spans tracer.Spans) []*sentry.Span {
	var result []*sentry.Span
	for _, span := range spans {
		if tracer.IsNoopSpan(span) {
			continue
		}
		entry := &sentry.Span{
			StartTime:   span.StartTS(),
			Description: span.Name(),
			Status:      sentry.SpanStatusOK,
			Data:        map[string]interface{}{},
		}
		span.Fields().ForEachField(func(f *field.Field) bool {
			entry.Data[f.Key] = f.Value
			return true
		})
		traceIDs := span.TraceIDs()
		if len(traceIDs) > 0 {
			copy(entry.TraceID[:], traceIDs[0])
		}
		copy(entry.SpanID[:], fmt.Sprint(span.ID()))
		if parent := span.Parent(); parent != nil {
			copy(entry.ParentSpanID[:], fmt.Sprint(parent.ID()))
		}
		result = append(result, entry)
	}
	return result
}

// UserToSentry converts an user structure to the Sentry format.
func UserToSentry(user *errmontypes.User) sentry.User {
	result := sentry.User{
		ID:       fmt.Sprint(user.ID),
		Username: user.Name,
	}

	for _, v := range user.CustomData {
		switch v := v.(type) {
		case UserEmail:
			result.Email = string(v)
		case UserIPAddress:
			result.IPAddress = string(v)
		}
	}

	return result
}

// HTTPRequestToSentry converts HTTP request info to the Sentry format.
func HTTPRequestToSentry(request *errmon.HTTPRequest) *sentry.Request {
	headers := make(map[string]string, len(request.Header))
	for name, values := range request.Header {
		headers[name] = strings.Join(values, "\n")
	}
	return &sentry.Request{
		URL:         request.URL.String(),
		Method:      request.Method,
		QueryString: request.URL.RawQuery,
		Cookies:     headers["Cookie"],
		Headers:     headers,
	}
}

// BreadcrumbToSentry converts a Breadcrumb to the Sentry format.
func BreadcrumbToSentry(breadcrumb *errmontypes.Breadcrumb) *sentry.Breadcrumb {
	data := map[string]interface{}{}
	breadcrumb.ForEachField(func(f *field.Field) bool {
		data[f.Key] = f.Value
		return true
	})
	return &sentry.Breadcrumb{
		Type:      strings.Join(breadcrumb.Path, "."),
		Category:  strings.Join(breadcrumb.Categories, ","),
		Data:      data,
		Timestamp: breadcrumb.TS,
	}
}

// PackageToSentry converts a Package to the Sentry format.
func PackageToSentry(pkg *errmontypes.Package) sentry.SdkPackage {
	return sentry.SdkPackage{
		Name:    pkg.Name,
		Version: pkg.Version,
	}
}

// EventToSentry converts an Event to the Sentry format.
func EventToSentry(ev *errmontypes.Event) *sentry.Event {
	result := &sentry.Event{
		EventID:  sentry.EventID(ev.ID),
		Level:    LevelToSentry(ev.Level),
		Message:  ev.Message,
		Platform: "go",
		Sdk: sentry.SdkInfo{
			Name: "go-belt",
		},
		Threads:   GoroutinesToSentry(ev.Goroutines, ev.CurrentGoroutineID),
		Timestamp: ev.Timestamp,

		Tags:    map[string]string{},
		Modules: map[string]string{},
		Extra:   map[string]interface{}{},

		StartTime: ev.Spans.Earliest().StartTS(),
		Spans:     SpansToSentry(ev.Spans),
	}

	if ev.Error != nil {
		result.Exception = append(result.Exception, sentry.Exception{
			Type:       "error",
			Value:      ev.Error.Error(),
			Module:     FuncNameToSentryModule(ev.Caller.Func().Name()),
			ThreadID:   "goroutine",
			Stacktrace: StackTraceToSentry(ev.StackTrace),
		})
	}
	if ev.IsPanic {
		result.Exception = append(result.Exception, sentry.Exception{
			Type:       "panic",
			Value:      fmt.Sprint(ev.PanicValue),
			Module:     FuncNameToSentryModule(ev.Caller.Func().Name()),
			ThreadID:   "goroutine",
			Stacktrace: StackTraceToSentry(ev.StackTrace),
		})
	}

	observeField := func(f *field.Field) bool {
		switch value := f.Value.(type) {
		case errmontypes.User:
			result.User = UserToSentry(&value)
		case *errmontypes.User:
			result.User = UserToSentry(value)
		case errmontypes.HTTPRequest:
			result.Request = HTTPRequestToSentry(&value)
		case *errmontypes.HTTPRequest:
			result.Request = HTTPRequestToSentry(value)
		case errmontypes.Breadcrumb:
			result.Breadcrumbs = append(result.Breadcrumbs, BreadcrumbToSentry(&value))
		case *errmontypes.Breadcrumb:
			result.Breadcrumbs = append(result.Breadcrumbs, BreadcrumbToSentry(value))
		case errmontypes.Package:
			result.Sdk.Packages = append(result.Sdk.Packages, PackageToSentry(&value))
		case *errmontypes.Package:
			result.Sdk.Packages = append(result.Sdk.Packages, PackageToSentry(value))
		case errmontypes.Tag:
			result.Tags[value.Key] = value.Value
		case *errmontypes.Tag:
			result.Tags[value.Key] = value.Value
		case tracer.SpanOptionRole:
			result.Type = string(value)
		case *tracer.SpanOptionRole:
			result.Type = string(*value)
		default:
			switch {
			case f.Properties.Has(errmontypes.FieldPropEnvironment):
				result.Environment = fmt.Sprint(f.Value)
			case f.Properties.Has(errmontypes.FieldPropRelease):
				result.Release = fmt.Sprint(f.Value)
			case f.Properties.Has(errmontypes.FieldPropServerName):
				result.ServerName = fmt.Sprint(f.Value)
			default:
				result.Extra[f.Key] = f.Value
			}
		}

		return true
	}

	ev.Fields.ForEachField(observeField)
	return result
}
