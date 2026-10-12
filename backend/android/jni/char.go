//go:build android

package jni

// Char builds a Java char argument, which is one UTF-16 unit and not a rune.
func Char(v uint16) Value { return Value(v) }
