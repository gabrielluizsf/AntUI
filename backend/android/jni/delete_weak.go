//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// DeleteWeak gives a weak reference back.
func (e *Env) DeleteWeak(o Object) {
	if o.ref != 0 {
		C.mw_DeleteWeakGlobalRef(e.ptr, o.ref)
	}
}
