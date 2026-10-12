//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid
#include <android/native_window.h>
*/
import "C"

// Release drops a reference taken by Acquire.
func (w *Window) Release() { C.ANativeWindow_release(w.ptr) }
