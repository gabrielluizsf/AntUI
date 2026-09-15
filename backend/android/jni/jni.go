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
	"sync"
	"unsafe"
)

// The virtual machine and the activity, set once by the app package when
// Android creates the activity. They are the only globals: everything else
// is per-thread or per-call by nature.
var (
	stateMu  sync.RWMutex
	vm       *C.JavaVM
	activity C.jobject // a global reference, so it outlives onCreate
)

// ErrNoVM is what every call reports before [Init] has run — which is every
// call in a process that Android did not start as an app.
var ErrNoVM = errors.New("jni: no Java virtual machine; this process is not an Android app")

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

// Ready reports whether there is a virtual machine to talk to.
func Ready() bool {
	stateMu.RLock()
	defer stateMu.RUnlock()
	return vm != nil
}

// Activity is the app's android.app.Activity, as a global reference. It is
// the starting point for almost everything: a Context, a window, a place to
// ask for a system service.
func Activity() Object {
	stateMu.RLock()
	defer stateMu.RUnlock()
	return Object{ref: activity}
}

func machine() *C.JavaVM {
	stateMu.RLock()
	defer stateMu.RUnlock()
	return vm
}

// Env is one thread's connection to the virtual machine, inside one local
// reference frame. It is only valid for the duration of the [Do] that made
// it — keeping one and using it later is a crash, not an error.
type Env struct {
	ptr *C.JNIEnv
}

// frameSize is how many local references a call may make before it has to
// manage them itself. The default table is 512 and overflowing it aborts the
// process, so a frame is opened for every Do; 32 is enough for anything that
// is not looping over an array, and the ones that do open their own.
const frameSize = 32

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

// attach finds this thread's connection, joining it to the machine if it is
// not already, and reports whether it has to be detached again.
func attach(m *C.JavaVM) (env *C.JNIEnv, detach bool, err error) {
	switch C.mw_GetEnv(m, &env, C.JNI_VERSION_1_6) {
	case C.JNI_OK:
		return env, false, nil
	case C.JNI_EDETACHED:
		if C.mw_AttachCurrentThread(m, &env) != C.JNI_OK {
			return nil, false, errors.New("jni: cannot attach this thread to the virtual machine")
		}
		return env, true, nil
	}
	return nil, false, errors.New("jni: the virtual machine does not speak JNI 1.6")
}

// Hold attaches the calling goroutine's thread and keeps it attached until
// the returned function is called. [Do] inside the hold is then only a
// lookup rather than an attach, which is what makes it cheap enough to use
// in a frame loop.
//
// It is also how a **sequence** of calls is kept on one thread, which some
// of Java's API quietly requires: a SQLite transaction belongs to the thread
// that began it, and a lock taken on one thread cannot be released on
// another. Go moves a goroutine between threads whenever it likes, and each
// [Do] only pins for its own duration — so anything that must happen on one
// thread from beginning to end has to hold.
//
// The goroutine stays pinned to its thread for the whole hold. **The
// returned function must be called before the goroutine ends**: a thread
// that dies while still attached takes the whole process down with it, and
// the message the runtime prints does not say so.
func Hold() (release func(), err error) {
	m := machine()
	if m == nil {
		return nil, ErrNoVM
	}
	runtime.LockOSThread()
	_, detach, err := attach(m)
	if err != nil {
		runtime.UnlockOSThread()
		return nil, err
	}
	var once sync.Once
	return func() {
		once.Do(func() {
			if detach {
				C.mw_DetachCurrentThread(m)
			}
			runtime.UnlockOSThread()
		})
	}, nil
}

// clearPending drops an exception without looking at it. It is only for the
// places where looking would need another call and that call is exactly what
// cannot be made.
func clearPending(env *C.JNIEnv) {
	if C.mw_ExceptionCheck(env) != C.JNI_FALSE {
		C.mw_ExceptionClear(env)
	}
}
