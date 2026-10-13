//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

import "errors"

// UnregisterNatives takes every native method off a class. The platform's
// own documentation calls it a tool for debugging and tool builders, and
// says an app should not need it; it is here for symmetry.
func (e *Env) UnregisterNatives(c Class) error {
	if C.mw_UnregisterNatives(e.ptr, c.ref) != C.JNI_OK {
		if err := e.check(); err != nil {
			return err
		}
		return errors.New("jni: unregistering natives failed")
	}
	return nil
}
