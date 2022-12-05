package types

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEventGetID(t *testing.T) {
	(*Event)(nil).GetID() // should not panic
}

func TestEventError(t *testing.T) {
	ev := &Event{
		ID: "1",
		Exception: Exception{
			StackTrace: []runtime.Frame{
				{
					Function: "func1()",
					File:     "file1.go",
					Line:     123,
				},
				{
					Function: "func2()",
					File:     "file2.go",
					Line:     234,
				},
			},
		},
		CurrentGoroutineID: 2,
		Goroutines: []Goroutine{
			{
				ID: 2,
			},
			{
				ID: 3,
			},
		},
	}

	t.Run("panic", func(t *testing.T) {
		ev.Exception.Error = nil
		ev.Exception.IsPanic = true
		ev.Exception.PanicValue = "unit-test"
		err := ev.AsError()
		require.ErrorAs(t, err, &ErrPanic{})
		require.Equal(t, `[event ID: '1'][goroutine 2] got a panic: unit-test
stack trace:
1. file1.go:123: func1()
2. file2.go:234: func2()
`, err.Error())
	})

	t.Run("error", func(t *testing.T) {
		ev.Exception.Error = fmt.Errorf("unit-test")
		ev.Exception.IsPanic = false
		ev.Exception.PanicValue = nil
		err := ev.AsError()
		require.ErrorAs(t, err, &ErrError{})
		require.Equal(t, "[event ID: '1'] unit-test", err.Error())
	})
}
