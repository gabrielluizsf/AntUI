//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid
#include <android/native_window.h>
*/
import "C"

// Width and Height are the size in pixels of the buffers this window hands
// out — which is the size asked for in SetGeometry when one was asked for,
// and the physical size of the surface otherwise.
func (w *Window) Width() int { return int(C.ANativeWindow_getWidth(w.ptr)) }
