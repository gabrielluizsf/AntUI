//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// CallBool calls a method that returns a Java boolean.
func (e *Env) CallBool(o Object, m Method, args ...Value) (bool, error) {
	if o.IsNil() {
		return false, errNilObject
	}
	a := jvalues(args)
	r := C.mw_CallBoolean(e.ptr, o.ref, m.id, first(a))
	keep(a)
	return r != C.JNI_FALSE, e.check()
}
