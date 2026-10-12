//go:build android

package jni

// Handle is the raw reference, for handing to code that speaks JNI itself.
func (o Object) Handle() uintptr { return uintptr(o.ref) }
