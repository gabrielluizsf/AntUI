//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// CallStaticLong calls a class method that returns a Java long.
func (e *Env) CallStaticLong(c Class, m Method, args ...Value) (int64, error) {
	a := jvalues(args)
	r := C.mw_CallStaticLong(e.ptr, c.ref, m.id, first(a))
	keep(a)
	return int64(r), e.check()
}
