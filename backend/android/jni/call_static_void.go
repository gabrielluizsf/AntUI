//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// CallStaticVoid calls a class method that returns nothing.
func (e *Env) CallStaticVoid(c Class, m Method, args ...Value) error {
	a := jvalues(args)
	C.mw_CallStaticVoid(e.ptr, c.ref, m.id, first(a))
	keep(a)
	return e.check()
}
