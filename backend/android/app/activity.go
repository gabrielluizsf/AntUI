//go:build android

package app

/*
#cgo LDFLAGS: -landroid -llog

#include <android/native_activity.h>
#include "glue.h"
*/
import "C"

import (
	"runtime"
	"runtime/debug"
	"sync"
	"unsafe"

	"github.com/gabrielluizsf/antui/backend/android/jni"
	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

// There is exactly one activity in a process, so it is a package variable
// rather than something threaded through every callback. It is written on
// the UI thread in onCreate and onDestroy and read there and on the app's
// goroutine, so it is behind a lock.
var (
	stateMu  sync.Mutex
	theApp   *App
	activity *C.ANativeActivity
)

func current() *App {
	stateMu.Lock()
	defer stateMu.Unlock()
	return theApp
}

// ANativeActivity_onCreate is the one function Android looks for in this
// library. It runs on the UI thread and has to return quickly, so all it
// does is set the callbacks up and start the goroutine that is the app.
//
//export ANativeActivity_onCreate
func ANativeActivity_onCreate(act *C.ANativeActivity, savedState unsafe.Pointer, savedStateSize C.size_t) {
	// Before anything else, so that a panic in the lines below is readable
	// in logcat instead of vanishing.
	ndk.RedirectStdio()

	// This is the UI thread, and the only place its looper can be reached.
	startUIRunner()

	f := registeredMain()
	if f == nil {
		ndk.Errorf("no app registered: call app.Main(f) from an init function")
		return
	}

	a := &App{
		// A little buffering, so a burst of configuration changes does not
		// block the UI thread on an app that is mid-frame.
		msgs: make(chan message, 16),
		dead: make(chan struct{}),
		info: Info{
			InternalData: C.GoString(act.internalDataPath),
			ExternalData: C.GoString(act.externalDataPath),
			OBB:          C.GoString(act.obbPath),
			SDK:          int(act.sdkVersion),
		},
		config: readConfig(act),
		assets: unsafe.Pointer(act.assetManager),
		state:  Created,
		saved:  copyState(savedState, savedStateSize),
	}

	stateMu.Lock()
	theApp, activity = a, act
	stateMu.Unlock()

	// The bridge to Java, before anything else can want it. This runs on the
	// UI thread, which is the only place the activity reference handed over
	// by the platform is valid — the first thing Init does is make it
	// global, and everything else in the library reaches Java through that.
	if err := jni.Init(unsafe.Pointer(act.vm), uintptr(act.clazz)); err != nil {
		ndk.Errorf("no bridge to Java: %v", err)
	} else {
		registerShim()
		// Before the app's goroutine starts, so that the first thing it does
		// can already make a temporary file.
		setupPaths(a)
	}

	C.antui_set_callbacks(act)

	go func() {
		// The looper, the input queue and most of the NDK are thread-local,
		// so the app's goroutine keeps one thread for its whole life.
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		run(a, f)
	}()
}

// copyState takes a copy of what a previous run saved. The block belongs to
// the framework and is freed when onCreate returns, so it cannot simply be
// pointed at.
func copyState(p unsafe.Pointer, n C.size_t) []byte {
	if p == nil || n == 0 {
		return nil
	}
	out := make([]byte, int(n))
	copy(out, unsafe.Slice((*byte)(p), int(n)))
	return out
}

// run calls the app's function and then keeps the machinery answering.
func run(a *App, f func(*App)) {
	// Keep this thread attached to the virtual machine for as long as the
	// app runs. Attaching allocates a thread object inside the runtime, and
	// doing it per call would put that cost in the middle of a frame loop.
	if release, err := jni.Hold(); err == nil {
		defer release()
	}
	defer func() {
		if r := recover(); r != nil {
			ndk.Errorf("panic: %v\n%s", r, debug.Stack())
		}
		// Whether it returned or panicked, the app is over: ask Android to
		// close the activity, then go on acknowledging callbacks until it
		// does. Skipping this is how an app hangs instead of exiting.
		Finish()
		a.drain()
	}()
	f(a)
}

// Finish asks Android to close the activity. It returns immediately; the
// Destroy event arrives later, in the usual order after Pause and Stop.
func Finish() {
	stateMu.Lock()
	act := activity
	stateMu.Unlock()
	if act != nil {
		C.ANativeActivity_finish(act)
	}
}

// post hands an event to the app's goroutine. When wait is set it does not
// return until the app has acknowledged it, which is what makes it safe for
// Android to take the window away on the next line.
func post(ev Event, ptr unsafe.Pointer, wait bool) {
	a := current()
	if a == nil {
		return
	}
	m := message{ev: ev, ptr: ptr}
	if wait {
		m.ack = make(chan struct{})
	}
	select {
	case a.msgs <- m:
	case <-a.dead:
		return
	}
	if m.ack == nil {
		return
	}
	select {
	case <-m.ack:
	case <-a.dead:
	}
}

// The lifecycle. Each of these blocks the UI thread until the app has seen
// it: an app that misses its Pause does not get another chance, because the
// process may be frozen the moment the callback returns.

//export antuiOnStart
func antuiOnStart(*C.ANativeActivity) { post(Start, nil, true) }

//export antuiOnResume
func antuiOnResume(*C.ANativeActivity) { post(Resume, nil, true) }

//export antuiOnPause
func antuiOnPause(*C.ANativeActivity) { post(Pause, nil, true) }

//export antuiOnStop
func antuiOnStop(*C.ANativeActivity) { post(Stop, nil, true) }

//export antuiOnDestroy
func antuiOnDestroy(*C.ANativeActivity) {
	post(Destroy, nil, true)
	stateMu.Lock()
	theApp, activity = nil, nil
	stateMu.Unlock()
}

// The window. Created and destroyed both wait, and destroyed is the one that
// matters: the surface is freed as soon as this returns, so the app must
// have stopped drawing into it first.

//export antuiOnNativeWindowCreated
func antuiOnNativeWindowCreated(_ *C.ANativeActivity, w *C.ANativeWindow) {
	post(WindowUp, unsafe.Pointer(w), true)
}

//export antuiOnNativeWindowResized
func antuiOnNativeWindowResized(_ *C.ANativeActivity, w *C.ANativeWindow) {
	post(WindowResized, unsafe.Pointer(w), false)
}

//export antuiOnNativeWindowRedrawNeeded
func antuiOnNativeWindowRedrawNeeded(_ *C.ANativeActivity, w *C.ANativeWindow) {
	// This one does not wait, and it is the one place that is a deliberate
	// loss. Waiting until the frame is actually on screen is what stops a
	// rotation flashing black — but the frame is drawn by the app between
	// pump and present, and holding the UI thread across that is holding it
	// across arbitrary user code, which may perfectly reasonably ask the UI
	// thread for the window insets and then wait for a thread that is
	// waiting for it. A flash is a blemish; a deadlock is the app gone.
	post(Redraw, unsafe.Pointer(w), false)
}

//export antuiOnNativeWindowDestroyed
func antuiOnNativeWindowDestroyed(_ *C.ANativeActivity, w *C.ANativeWindow) {
	post(WindowDown, unsafe.Pointer(w), true)
}

// The input queue. Reading it belongs to the input layer; all this does is
// say when it exists.

//export antuiOnInputQueueCreated
func antuiOnInputQueueCreated(_ *C.ANativeActivity, q *C.AInputQueue) {
	post(InputUp, unsafe.Pointer(q), true)
}

//export antuiOnInputQueueDestroyed
func antuiOnInputQueueDestroyed(_ *C.ANativeActivity, q *C.AInputQueue) {
	post(InputDown, unsafe.Pointer(q), true)
}

// The rest is news, not instruction: none of it blocks.

//export antuiOnWindowFocusChanged
func antuiOnWindowFocusChanged(_ *C.ANativeActivity, hasFocus C.int) {
	if hasFocus != 0 {
		post(FocusGained, nil, false)
	} else {
		post(FocusLost, nil, false)
	}
}

//export antuiOnContentRectChanged
func antuiOnContentRectChanged(_ *C.ANativeActivity, r *C.ARect) {
	if a := current(); a != nil && r != nil {
		a.mu.Lock()
		a.rect = ndk.Rect{
			Left:   int32(r.left),
			Top:    int32(r.top),
			Right:  int32(r.right),
			Bottom: int32(r.bottom),
		}
		a.mu.Unlock()
	}
	post(ContentRect, nil, false)
}

//export antuiOnConfigurationChanged
func antuiOnConfigurationChanged(act *C.ANativeActivity) {
	// A configuration change is where the density moves: a window dragged to
	// another display, or a device that reports differently once it is
	// unfolded.
	if a := current(); a != nil {
		cfg := readConfig(act)
		a.mu.Lock()
		a.config = cfg
		a.mu.Unlock()
	}
	post(ConfigChanged, nil, false)
}

//export antuiOnLowMemory
func antuiOnLowMemory(*C.ANativeActivity) { post(LowMemory, nil, false) }
