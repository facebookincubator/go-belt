package types

import (
	"fmt"
)

// ErrPanic is a wrapper of an Event which implements `error` in case if
// the Event conveys a panic event.
type ErrPanic struct {
	Event *Event
}

// Error implements interface `error`.
func (err ErrPanic) Error() string {
	return fmt.Sprintf(`[event ID: '%s'][goroutine %d] got a panic: %v
stack trace:
%s`,
		err.Event.ID,
		err.Event.CurrentGoroutineID,
		err.Event.PanicValue,
		err.Event.StackTrace,
	)
}

// ErrError is a wrapper of an Event which implements `error` in case if
// the Event conveys an error event.
type ErrError struct {
	Event *Event
}

// Error implements interface `error`.
func (err ErrError) Error() string {
	return fmt.Sprintf("[event ID: '%s'] %v", err.Event.ID, err.Unwrap())
}

// Unwrap implements go1.13 errors unwrapping.
func (err ErrError) Unwrap() error {
	return err.Event.Exception.Error
}
