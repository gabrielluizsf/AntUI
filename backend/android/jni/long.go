//go:build android

package jni

// Long builds a Java long argument.
func Long(v int64) Value { return Value(uint64(v)) }
