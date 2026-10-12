//go:build android

package jni

import (
	"math"
)

// Float builds a Java float argument.
func Float(v float32) Value { return Value(math.Float32bits(v)) }
