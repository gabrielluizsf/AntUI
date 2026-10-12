//go:build android

package jni

// Int builds a Java int argument.
func Int(v int32) Value { return Value(uint32(v)) }
