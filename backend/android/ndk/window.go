//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid
#include <android/native_window.h>
*/
import "C"

// Window is an ANativeWindow: the surface the system gives an app to draw
// into. It is only ever created by Android and handed to a callback; this
// package never makes one.
//
// A Window is valid only between the system saying it exists and the system
// saying it does not. Drawing into one after that is a crash, not an error,
// which is why [Window.Acquire] exists.
type Window struct {
	ptr *C.ANativeWindow
}
