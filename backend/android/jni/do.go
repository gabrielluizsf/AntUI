//go:build android

package jni

/*
#cgo LDFLAGS: -landroid -llog

#include "jni_wrap.h"
*/
import "C"

import (
	"errors"
	"runtime"
)

// Do runs f with a connection to the virtual machine.
//
// It pins the goroutine to its thread for the whole call — a JNIEnv belongs
// to a thread and Go is otherwise free to move a goroutine between them,
// which would not fail, it would corrupt. If the thread is not attached to
// the machine it is attached and detached again; a thread that will make
// many calls should be held attached with [Hold] instead, because attaching
// allocates a thread object inside the runtime every time.
func Do(f func(*Env) error) error {
	m := machine()
	if m == nil {
		return ErrNoVM
	}
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	env, detach, err := attach(m)
	if err != nil {
		return err
	}
	if detach {
		defer C.mw_DetachCurrentThread(m)
	}

	if C.mw_PushLocalFrame(env, frameSize) != 0 {
		clearPending(env)
		return errors.New("jni: cannot open a local reference frame")
	}
	defer C.mw_PopLocalFrame(env, 0)

	return f(&Env{ptr: env})
}
