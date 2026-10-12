//go:build android

package app

import (
	"runtime/debug"

	"github.com/gabrielluizsf/antui/backend/android/jni"
	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

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
