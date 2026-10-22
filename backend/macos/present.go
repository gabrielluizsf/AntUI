//go:build darwin

package macos

/*
#cgo CFLAGS: -x objective-c -Wno-deprecated-declarations
#cgo LDFLAGS: -framework Cocoa -framework CoreGraphics -framework CoreVideo
#include <stdlib.h>
#include "backend_darwin.h"
*/
import "C"
import (
	"unsafe"

	"github.com/gabrielluizsf/antui/backend"
	"github.com/gabrielluizsf/antui/canvas"
)

// Present copies the window's dirty rectangle to the platform surface,
// reporting whether there was a surface to copy it to.
func (d *Driver) Present(win backend.Face, dirty canvas.Area) bool {
	if d.handle == nil || dirty.Width <= 0 || dirty.Height <= 0 {
		return false
	}
	cv := win.Canvas()
	if cv.Pixels == nil {
		return false
	}
	// The pixels are copied on the other side, not retained, so handing a Go
	// pointer across for the length of this call is within the cgo rules.
	C.antui_d_present(d.handle,
		(*C.uint32_t)(unsafe.Pointer(&cv.Pixels[0])),
		C.int(cv.Width), C.int(cv.Height),
		C.int(dirty.X), C.int(dirty.Y), C.int(dirty.Width), C.int(dirty.Height))
	return true
}
