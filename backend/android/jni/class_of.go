//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

import (
	"errors"
)

// ClassOf is the class of an object, which is how a returned object's real
// type is found without knowing its name.
//
// **The caller owns what comes back and has to let it go**:
//
//	cls, err := e.ClassOf(o)
//	if err != nil {
//	    return err
//	}
//	defer e.Delete(cls.Object())
//
// It is a *local* reference, unlike the one [Env.Class] hands out, and two
// things follow from that. It has to be deleted, or a loop fills the local
// reference table — which is 512 entries deep and then the VM aborts. And it
// must never be used as the key of anything that outlives the call: the slot
// is reused, so the same number means one class now and a different class in
// a minute. Both of those were live bugs here, and neither showed up in a
// year of running until CheckJNI was turned on.
func (e *Env) ClassOf(o Object) (Class, error) {
	if o.IsNil() {
		return Class{}, errors.New("jni: no class for a null object")
	}
	ref := C.mw_GetObjectClass(e.ptr, o.ref)
	if ref == 0 {
		return Class{}, e.check()
	}
	return Class{ref: ref}, nil
}
