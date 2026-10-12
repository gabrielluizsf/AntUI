//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// CallStaticFloat calls a class method that returns a Java float.
func (e *Env) CallStaticFloat(c Class, m Method, args ...Value) (float32, error) {
	a := jvalues(args)
	r := C.mw_CallStaticFloat(e.ptr, c.ref, m.id, first(a))
	keep(a)
	return float32(r), e.check()
}
