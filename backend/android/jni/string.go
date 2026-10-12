//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

import (
	"runtime"
	"unicode/utf16"
	"unsafe"
)

// String makes a Java String out of a Go one.
//
// It goes through UTF-16 and not through NewStringUTF, which is the obvious
// call and the wrong one. NewStringUTF takes **modified UTF-8**, which is
// not UTF-8: a NUL is encoded as two bytes rather than one, and anything
// outside the basic plane — every emoji — is encoded as a surrogate pair of
// three-byte sequences rather than as one four-byte sequence. Handing it a
// Go string containing either produces a Java String that is quietly wrong,
// or a crash inside the runtime. Java's own representation is UTF-16, so
// converting to that is both correct and what it was going to do anyway.
func (e *Env) String(s string) (Object, error) {
	u := utf16.Encode([]rune(s))
	if len(u) == 0 {
		r := C.mw_NewString(e.ptr, nil, 0)
		return Object{ref: C.jobject(r)}, e.check()
	}
	r := C.mw_NewString(e.ptr, (*C.jchar)(unsafe.Pointer(&u[0])), C.jsize(len(u)))
	runtime.KeepAlive(u)
	return Object{ref: C.jobject(r)}, e.check()
}
