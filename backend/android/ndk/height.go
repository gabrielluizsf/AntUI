//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid
#include <android/native_window.h>
*/
import "C"

// Height is the surface in pixels, top to bottom.
func (w *Window) Height() int { return int(C.ANativeWindow_getHeight(w.ptr)) }
