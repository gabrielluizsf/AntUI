//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid
#include <android/native_window.h>
*/
import "C"

import (
	"unsafe"
)

// Lock takes the next buffer to draw into. When dirty is not nil it names
// the region about to be written, and Lock rewrites it with the region the
// compositor actually requires — which is usually larger, because it has to
// include whatever the previous frames in the swap chain left behind. Code
// that trusts its own rectangle instead of the one that comes back leaves
// stale pixels on screen every second or third frame.
func (w *Window) Lock(dirty *Rect) (Buffer, error) {
	var buf C.ANativeWindow_Buffer
	var rect *C.ARect
	if dirty != nil {
		rect = (*C.ARect)(unsafe.Pointer(dirty))
	}
	if rc := C.ANativeWindow_lock(w.ptr, &buf, rect); rc != 0 {
		return Buffer{}, wrapErr("ANativeWindow_lock", int32(rc))
	}
	return Buffer{
		Width:  int(buf.width),
		Height: int(buf.height),
		Stride: int(buf.stride),
		Format: Format(buf.format),
		Bits:   buf.bits,
	}, nil
}
