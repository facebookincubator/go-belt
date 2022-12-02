package types

import (
	"testing"
)

func TestEventGetID(t *testing.T) {
	(*Event)(nil).GetID() // should not panic
}
