package assert

import "testing"

func Error(t *testing.T, err error) {
	if err == nil {
		t.Fatal("want error")
	}
}
