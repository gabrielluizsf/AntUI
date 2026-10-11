//go:build darwin

package macos

/*
#cgo CFLAGS: -x objective-c -Wno-deprecated-declarations
#cgo LDFLAGS: -framework Cocoa -framework CoreGraphics -framework CoreVideo
#include <stdlib.h>
#include "backend_darwin.h"
*/
import "C"
import "github.com/gabrielluizsf/antui/backend"

// SetSize asks the platform for a new drawable size, reporting whether the
// request could be made.
func (d *Driver) SetSize(win backend.Face, width, height int) bool {
	if d.handle == nil {
		return false
	}
	return C.antui_d_set_size(d.handle, C.int(width), C.int(height)) != 0
}
