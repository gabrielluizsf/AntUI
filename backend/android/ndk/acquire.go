//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid
#include <android/native_window.h>
*/
import "C"

// Acquire adds a reference, so the window outlives the callback that
// announced it. Every Acquire needs a Release.
func (w *Window) Acquire() { C.ANativeWindow_acquire(w.ptr) }
