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

	"github.com/gabrielluizsf/antui/canvas"
)

// SetIcon gives the window an icon where the system shows one. The largest
// picture is handed to AppKit, which is the one the Dock wants: macOS has no
// small icon of its own to set.
func (d *Driver) SetIcon(images []*canvas.Canvas) bool {
	if d.handle == nil || len(images) == 0 {
		return false
	}
	best := images[0]
	for _, image := range images {
		if image.Width > best.Width {
			best = image
		}
	}
	pixels := make([]uint32, best.Width*best.Height)
	for y := range best.Height {
		for x := range best.Width {
			pixels[y*best.Width+x] = uint32(best.At(x, y))
		}
	}
	C.antui_d_set_icon(d.handle, (*C.uint32_t)(unsafe.Pointer(&pixels[0])),
		C.int(best.Width), C.int(best.Height))
	return true
}
