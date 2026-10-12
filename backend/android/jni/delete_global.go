//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// DeleteGlobal gives a global reference back.
func (e *Env) DeleteGlobal(o Object) {
	if o.ref != 0 {
		C.mw_DeleteGlobalRef(e.ptr, o.ref)
	}
}
