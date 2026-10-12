//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// CallDouble calls a method that returns a Java double.
func (e *Env) CallDouble(o Object, m Method, args ...Value) (float64, error) {
	if o.IsNil() {
		return 0, errNilObject
	}
	a := jvalues(args)
	r := C.mw_CallDouble(e.ptr, o.ref, m.id, first(a))
	keep(a)
	return float64(r), e.check()
}
