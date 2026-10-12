//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// CallStaticDouble calls a class method that returns a Java double.
func (e *Env) CallStaticDouble(c Class, m Method, args ...Value) (float64, error) {
	a := jvalues(args)
	r := C.mw_CallStaticDouble(e.ptr, c.ref, m.id, first(a))
	keep(a)
	return float64(r), e.check()
}
