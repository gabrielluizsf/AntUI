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

// GoString reads a Java String.
//
// A null String comes back as the empty string with no error, because that
// is what almost every caller wants and the alternative is a second return
// value on every read. Where the difference matters, [Object.IsNil] says so.
func (e *Env) GoString(o Object) (string, error) {
	if o.IsNil() {
		return "", nil
	}
	n := C.mw_GetStringLength(e.ptr, C.jstring(o.ref))
	if err := e.check(); err != nil {
		return "", err
	}
	if n == 0 {
		return "", nil
	}
	buf := make([]uint16, int(n))
	C.mw_GetStringRegion(e.ptr, C.jstring(o.ref), 0, n, (*C.jchar)(unsafe.Pointer(&buf[0])))
	if err := e.check(); err != nil {
		return "", err
	}
	// Decode pairs the surrogates back up. A lone surrogate — which a Java
	// String is allowed to contain and a Go string is not — becomes U+FFFD
	// rather than being carried through as something no Go code can handle.
	return string(utf16.Decode(buf)), nil
}
