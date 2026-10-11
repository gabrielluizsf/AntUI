//go:build darwin

package macos

/*
#cgo CFLAGS: -x objective-c -Wno-deprecated-declarations
#cgo LDFLAGS: -framework Cocoa -framework CoreGraphics -framework CoreVideo
#include <stdlib.h>
#include "backend_darwin.h"
*/
import "C"

// DisplayRefresh is how often the display repaints, in hertz.
func (d *Driver) DisplayRefresh() int {
	return int(C.antui_d_display_refresh(d.handle))
}
