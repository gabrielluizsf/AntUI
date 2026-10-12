//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid
#include <android/native_window.h>
*/
import "C"

// UnlockAndPost hands the buffer back and puts it on screen.
func (w *Window) UnlockAndPost() error {
	return wrapErr("ANativeWindow_unlockAndPost", int32(C.ANativeWindow_unlockAndPost(w.ptr)))
}
