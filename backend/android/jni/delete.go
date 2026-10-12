//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// Delete frees one local reference early. It is only worth doing inside a
// loop that makes many — a frame frees the rest.
func (e *Env) Delete(o Object) {
	if o.ref != 0 {
		C.mw_DeleteLocalRef(e.ptr, o.ref)
	}
}
