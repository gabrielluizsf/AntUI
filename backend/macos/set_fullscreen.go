//go:build darwin

package macos

/*
#cgo CFLAGS: -x objective-c -Wno-deprecated-declarations
#cgo LDFLAGS: -framework Cocoa -framework CoreGraphics -framework CoreVideo
#include <stdlib.h>
#include "backend_darwin.h"
*/
import "C"

// SetFullscreen asks the system for a full-screen window, reporting whether
// the request could be made at all.
func (d *Driver) SetFullscreen(on bool) bool {
	if d.handle == nil {
		return false
	}
	flag := C.int(0)
	if on {
		flag = 1
	}
	return C.antui_d_set_fullscreen(d.handle, flag) != 0
}
