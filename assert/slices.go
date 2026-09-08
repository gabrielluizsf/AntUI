package assert

import "testing"

func Len[T any, S []T](t *testing.T, expected int, s S) {
	if len(s) != expected {
		t.Fatalf("{ expected: %d, result: %d }", expected, len(s))
	}
}
