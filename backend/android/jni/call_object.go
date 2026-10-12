//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// CallObject calls a method that returns an object. The result is
// a local reference and lives until the frame it was made in ends — see
// [Env.NewGlobal] for one that has to outlive the call.
func (e *Env) CallObject(o Object, m Method, args ...Value) (Object, error) {
	if o.IsNil() {
		return Object{}, errNilObject
	}
	a := jvalues(args)
	r := C.mw_CallObject(e.ptr, o.ref, m.id, first(a))
	keep(a)
	return Object{ref: r}, e.check()
}
