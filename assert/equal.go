package assert

import (
	"reflect"
	"testing"
)

func Equal[T comparable](t *testing.T, expected, result T) {
	if expected != result {
		t.Fatalf("%v != %v", expected, result)
	}
}

func DeepEqual[T any](t *testing.T, expected, result T) {
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("%v != %v", expected, result)
	}
}