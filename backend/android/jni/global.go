//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// Global makes a reference that outlives the call and may be used from any
// thread. Every one has to be given back with [Env.DeleteGlobal]; they do
// not go away on their own and the table they live in is not unlimited.
func (e *Env) Global(o Object) Object {
	if o.ref == 0 {
		return Object{}
	}
	return Object{ref: C.mw_NewGlobalRef(e.ptr, o.ref)}
}
