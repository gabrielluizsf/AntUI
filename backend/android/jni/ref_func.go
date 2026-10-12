//go:build android

package jni

// Ref passes an object, or Java's null when the object is nil.
func Ref(o Object) Value { return Value(o.ref) }
