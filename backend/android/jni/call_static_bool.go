//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// CallStaticBool calls a class method that returns a Java boolean.
func (e *Env) CallStaticBool(c Class, m Method, args ...Value) (bool, error) {
	a := jvalues(args)
	r := C.mw_CallStaticBoolean(e.ptr, c.ref, m.id, first(a))
	keep(a)
	return r != C.JNI_FALSE, e.check()
}
