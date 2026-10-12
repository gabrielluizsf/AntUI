//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// CallStaticInt calls a class method that returns a Java int.
func (e *Env) CallStaticInt(c Class, m Method, args ...Value) (int32, error) {
	a := jvalues(args)
	r := C.mw_CallStaticInt(e.ptr, c.ref, m.id, first(a))
	keep(a)
	return int32(r), e.check()
}
