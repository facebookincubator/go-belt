package types

import (
	"fmt"
	"runtime"
	"strings"
)

// StackTrace is a goroutine execution frames stack trace.
type StackTrace []runtime.Frame

// String implements fmt.Stringer.
func (s StackTrace) String() string {
	var result strings.Builder
	for idx, frame := range s {
		result.WriteString(fmt.Sprintf("%d. %s:%d: %s\n", idx+1, frame.File, frame.Line, frame.Function))
	}
	return result.String()
}
