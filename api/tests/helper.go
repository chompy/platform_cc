package tests

import "testing"

func assertEqual(actual interface{}, expected interface{}, msg string, t *testing.T) {
	if actual != expected {
		t.Errorf("%s, got: %s, want: %s", msg, actual, expected)
	}
}
