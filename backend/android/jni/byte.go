//go:build android

package jni

// Byte builds a Java byte argument.
func Byte(v int8) Value { return Value(uint8(v)) }
