//go:build android

package jni

import (
	"math"
)

// Double builds a Java double argument.
func Double(v float64) Value {
	return Value(math.Float64bits(v))
}
