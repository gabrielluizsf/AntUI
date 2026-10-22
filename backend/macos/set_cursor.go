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
)

// SetCursor replaces the mouse pointer with icon's picture while it is over
// the window's content, reporting whether the request could be made. A nil
// Image puts the ordinary system arrow back.
func (d *Driver) SetCursor(icon backend.CursorIcon) bool {
	if d.handle == nil {
		return false
	}
	if icon.Image == nil || icon.Image.Width <= 0 || icon.Image.Height <= 0 {
		C.antui_d_set_cursor(d.handle, nil, 0, 0, 0, 0)
		return true
	}
	pixels := make([]uint32, icon.Image.Width*icon.Image.Height)
	for y := range icon.Image.Height {
		for x := range icon.Image.Width {
			pixels[y*icon.Image.Width+x] = uint32(icon.Image.At(x, y))
		}
	}
	C.antui_d_set_cursor(d.handle, (*C.uint32_t)(unsafe.Pointer(&pixels[0])),
		C.int(icon.Image.Width), C.int(icon.Image.Height),
		C.int(icon.HotX), C.int(icon.HotY))
	return true
}
