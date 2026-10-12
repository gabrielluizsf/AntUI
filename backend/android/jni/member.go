//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

import (
	"errors"
	"unsafe"
)

// member is the half every lookup shares: the C strings, the nil check, and
// an error that says what was being looked for. A lookup that fails leaves
// a NoSuchMethodError pending, which has to be cleared or the next call
// fails instead.
func (e *Env) member(c Class, name, sig string, find func(cn, cs *C.char) unsafe.Pointer, what string) (unsafe.Pointer, error) {
	if c.IsNil() {
		return nil, errors.New("jni: no " + what + " " + name + ": the class was not found")
	}
	cname, csig := C.CString(name), C.CString(sig)
	defer C.free(unsafe.Pointer(cname))
	defer C.free(unsafe.Pointer(csig))
	id := find(cname, csig)
	if id == nil {
		if err := e.check(); err != nil {
			return nil, errors.New("jni: no " + what + " " + name + sig + ": " + err.Error())
		}
		return nil, errors.New("jni: no " + what + " " + name + sig)
	}
	return id, nil
}
