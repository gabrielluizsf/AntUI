//go:build android

package jni

// Short builds a Java short argument.
func Short(v int16) Value { return Value(uint16(v)) }
