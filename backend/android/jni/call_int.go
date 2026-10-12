//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// CallInt calls a method that returns a Java int.
func (e *Env) CallInt(o Object, m Method, args ...Value) (int32, error) {
	if o.IsNil() {
		return 0, errNilObject
	}
	a := jvalues(args)
	r := C.mw_CallInt(e.ptr, o.ref, m.id, first(a))
	keep(a)
	return int32(r), e.check()
}
