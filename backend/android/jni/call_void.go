//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// CallVoid calls a method that returns nothing.
func (e *Env) CallVoid(o Object, m Method, args ...Value) error {
	if o.IsNil() {
		return errNilObject
	}
	a := jvalues(args)
	C.mw_CallVoid(e.ptr, o.ref, m.id, first(a))
	keep(a)
	return e.check()
}
