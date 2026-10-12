//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// CallLong calls a method that returns a Java long.
func (e *Env) CallLong(o Object, m Method, args ...Value) (int64, error) {
	if o.IsNil() {
		return 0, errNilObject
	}
	a := jvalues(args)
	r := C.mw_CallLong(e.ptr, o.ref, m.id, first(a))
	keep(a)
	return int64(r), e.check()
}
