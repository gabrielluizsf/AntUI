//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid
#include <android/native_window.h>
*/
import "C"

// SetGeometry asks for buffers of a given size and format. A width and
// height of zero means "whatever the surface is", which is what an app that
// draws at native resolution wants; giving a smaller size is how a game
// renders at half resolution and lets the compositor scale it up for free.
func (w *Window) SetGeometry(width, height int, f Format) error {
	rc := C.ANativeWindow_setBuffersGeometry(w.ptr, C.int32_t(width), C.int32_t(height), C.int32_t(f))
	return wrapErr("ANativeWindow_setBuffersGeometry", int32(rc))
}
