//go:build android

package jni

// IsNil reports whether there is no object — which is Java's null, and is
// what a method returns when it returns null.
func (o Object) IsNil() bool { return o.ref == 0 }
