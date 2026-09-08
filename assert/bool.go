package assert

import "testing"

func True(t *testing.T, condition bool) {
	if !condition {
		t.Fatal("{ expected: true, received: false }")
	}
}

func False(t *testing.T, condition bool) {
	if condition {
		t.Fatal("{ expected: false, received: true }")
	}
}
