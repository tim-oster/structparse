package structparse

import (
	"fmt"
	"testing"
)

func assertEqual(t *testing.T, expected, got interface{}) bool {
	var (
		expectedFmt = fmt.Sprintf("%v", expected)
		gotFmt      = fmt.Sprintf("%v", got)
	)
	if gotFmt != expectedFmt {
		t.Errorf("expected '%s' but got '%s'", expectedFmt, gotFmt)
		return false
	}
	return true
}

func Test_assertEqual(t *testing.T) {
	t.Run("unequal", func(t *testing.T) {
		dummy := &testing.T{}
		result := assertEqual(dummy, true, false)
		if result {
			t.Fatalf("wrong assert result: true != false")
		}
		if !dummy.Failed() {
			t.Fatalf("test should be marked as failed")
		}
	})

	t.Run("equal", func(t *testing.T) {
		dummy := &testing.T{}
		result := assertEqual(dummy, true, true)
		if !result {
			t.Fatalf("wrong assert result: true == true")
		}
		if dummy.Failed() {
			t.Fatalf("test should not be marked as failed")
		}
	})
}
