//go:build android

package jni

/*
#cgo LDFLAGS: -landroid -llog

#include "jni_wrap.h"
*/
import "C"

import (
	"errors"
	"unsafe"
)

// Init hands the bridge the virtual machine and the activity. The app
// package calls it from ANativeActivity_onCreate and nothing else should
// call it at all.
//
// The activity arrives as a local reference on the UI thread, which stops
// being valid the moment that callback returns, so the first thing done with
// it is to make it global.
// The activity is a uintptr and not an unsafe.Pointer because that is what
// it is: cgo models Java's opaque handle types — jobject and everything
// typedef'd from it — as uintptr rather than as pointers, since nothing on
// the Go side may ever dereference one.
func Init(javaVM unsafe.Pointer, activityObject uintptr) error {
	if javaVM == nil {
		return ErrNoVM
	}
	stateMu.Lock()
	vm = (*C.JavaVM)(javaVM)
	stateMu.Unlock()

	return Do(func(e *Env) error {
		ref := C.mw_NewGlobalRef(e.ptr, C.jobject(activityObject))
		if ref == 0 {
			return errors.New("jni: cannot hold on to the activity")
		}
		stateMu.Lock()
		activity = ref
		stateMu.Unlock()
		return nil
	})
}
