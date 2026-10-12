//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// Weak makes a reference that does not keep the object alive. It may read as
// nil at any time, so anything using one has to check.
func (e *Env) Weak(o Object) Object {
	if o.ref == 0 {
		return Object{}
	}
	return Object{ref: C.mw_NewWeakGlobalRef(e.ptr, o.ref)}
}
