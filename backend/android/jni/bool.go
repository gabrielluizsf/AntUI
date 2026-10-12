//go:build android

package jni

// The argument builders, one per Java type. A Value is one slot of a JNI
// argument list, and which builder made it is what says how the eight bytes
// are read on the other side — an int and a float of the same width are not
// the same argument.
func Bool(v bool) Value {
	if v {
		return 1
	}
	return 0
}
