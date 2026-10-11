//go:build darwin

package macos

/*
#cgo CFLAGS: -x objective-c -Wno-deprecated-declarations
#cgo LDFLAGS: -framework Cocoa -framework CoreGraphics -framework CoreVideo
#include <stdlib.h>
#include "backend_darwin.h"
*/
import "C"

// SetOpacity fades the whole window so what sits behind it shows through;
// 255 is fully opaque. AppKit applies the fade to the window as a whole, so
// the frame that was already drawn stays exactly as it was, only the window's
// overall brightness changes.
func (d *Driver) SetOpacity(alpha uint8) bool {
	if d.handle == nil {
		return false
	}
	C.antui_d_set_opacity(d.handle, C.int(alpha))
	return true
}
