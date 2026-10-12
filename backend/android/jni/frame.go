//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// Frame runs f inside a local reference frame of its own, so that everything
// f makes is freed when it returns. A loop that touches a thousand array
// elements needs this; anything else already has the frame [Do] opened.
func (e *Env) Frame(capacity int, f func() error) error {
	if C.mw_PushLocalFrame(e.ptr, C.jint(capacity)) != 0 {
		clearPending(e.ptr)
		return errNoFrame
	}
	err := f()
	C.mw_PopLocalFrame(e.ptr, 0)
	return err
}
