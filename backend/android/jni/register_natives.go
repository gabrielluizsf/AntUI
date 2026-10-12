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

// RegisterNatives attaches implementations to a class's native methods.
//
// A signature that does not match one the class declares is not an error
// here — the machine reports it as a NoSuchMethodError, which is what comes
// back — but a signature that matches the *wrong* method is neither, and is
// a crash the first time it is called.
func (e *Env) RegisterNatives(c Class, methods []Native) error {
	if len(methods) == 0 {
		return nil
	}
	if c.IsNil() {
		return errors.New("jni: cannot register natives on a class that was not found")
	}
	size := C.size_t(unsafe.Sizeof(C.JNINativeMethod{}))
	raw := C.malloc(size * C.size_t(len(methods)))
	if raw == nil {
		return errors.New("jni: out of memory registering natives")
	}
	defer C.free(raw)

	table := unsafe.Slice((*C.JNINativeMethod)(raw), len(methods))
	var owned []unsafe.Pointer
	defer func() {
		for _, p := range owned {
			C.free(p)
		}
	}()
	for i, m := range methods {
		name, sig := C.CString(m.Name), C.CString(m.Sig)
		owned = append(owned, unsafe.Pointer(name), unsafe.Pointer(sig))
		table[i].name = name
		table[i].signature = sig
		table[i].fnPtr = m.Fn
	}

	if C.mw_RegisterNatives(e.ptr, c.ref, (*C.JNINativeMethod)(raw), C.jint(len(methods))) != C.JNI_OK {
		if err := e.check(); err != nil {
			return err
		}
		return errors.New("jni: registering natives failed")
	}
	return nil
}
