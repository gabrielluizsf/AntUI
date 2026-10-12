//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// Wrap takes a raw jobject from somewhere else — the activity handed over by
// the platform, an argument to a native method. The caller keeps whatever
// lifetime it already had.
//
// The handle is a uintptr because that is how cgo models Java's opaque
// types. It is not an address Go may follow, and nothing here ever follows
// one.
func Wrap(p uintptr) Object { return Object{ref: C.jobject(p)} }
