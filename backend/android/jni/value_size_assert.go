//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

import "unsafe"

// A jvalue is a union containing a jlong and a jdouble, so it is eight bytes
// on every architecture Android has. Value being a uint64 depends on that,
// and these two lines are what would notice if it ever stopped being true.
var (
	_ [unsafe.Sizeof(C.jvalue{}) - 8]byte
	_ [8 - unsafe.Sizeof(C.jvalue{})]byte
)
