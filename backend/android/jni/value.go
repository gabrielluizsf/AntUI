//go:build android

package jni

// Value is one argument to a Java method: the eight bytes of a jvalue, which
// is a union of every type Java has.
//
// It is built with [Int], [Ref] and the rest rather than converted, because
// which eight bytes a number occupies depends on what it is. Passing an int
// where a long is expected is not a type error anywhere in JNI — it is a
// method that reads the wrong half of the stack.
type Value uint64
