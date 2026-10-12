//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// CallStaticObject calls a class method that returns an object.
func (e *Env) CallStaticObject(c Class, m Method, args ...Value) (Object, error) {
	a := jvalues(args)
	r := C.mw_CallStaticObject(e.ptr, c.ref, m.id, first(a))
	keep(a)
	return Object{ref: r}, e.check()
}
