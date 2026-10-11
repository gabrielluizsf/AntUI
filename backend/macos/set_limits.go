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

// SetLimits passes the size constraints on to the platform.
func (d *Driver) SetLimits(limits backend.Limits) {
	if d.handle == nil {
		return
	}
	fixed, noMaximize := C.int(0), C.int(0)
	if limits.Fixed {
		fixed = 1
	}
	if limits.NoMaximize {
		noMaximize = 1
	}
	C.antui_d_set_limits(d.handle,
		C.int(limits.MinWidth), C.int(limits.MinHeight),
		C.int(limits.MaxWidth), C.int(limits.MaxHeight),
		C.double(limits.Aspect), fixed, noMaximize)
}
