//go:build darwin

package macos

/*
#cgo CFLAGS: -x objective-c -Wno-deprecated-declarations
#cgo LDFLAGS: -framework Cocoa -framework CoreGraphics -framework CoreVideo
#include <stdlib.h>
#include "backend_darwin.h"
*/
import "C"

// DisplaySize is the size of the display the window is on.
func (d *Driver) DisplaySize() (w, h int, ok bool) {
	var width, height C.int
	if C.antui_d_display_size(d.handle, &width, &height) == 0 {
		return 0, 0, false
	}
	return int(width), int(height), true
}
