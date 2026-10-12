//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// CallFloat calls a method that returns a Java float.
func (e *Env) CallFloat(o Object, m Method, args ...Value) (float32, error) {
	if o.IsNil() {
		return 0, errNilObject
	}
	a := jvalues(args)
	r := C.mw_CallFloat(e.ptr, o.ref, m.id, first(a))
	keep(a)
	return float32(r), e.check()
}
