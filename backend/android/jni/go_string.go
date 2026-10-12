//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

import (
	"unicode/utf16"
	"unsafe"
)

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
