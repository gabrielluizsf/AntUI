//go:build darwin

package macos

/*
#cgo CFLAGS: -x objective-c -Wno-deprecated-declarations
#cgo LDFLAGS: -framework Cocoa -framework CoreGraphics -framework CoreVideo
#include <stdlib.h>
#include "backend_darwin.h"
*/
import "C"
import "runtime"

// Close destroys the window and lets go of everything the platform held.
func (d *Driver) Close() {
	if d.handle == nil {
		return
	}
	C.antui_d_close(d.handle)
	d.handle = nil
	runtime.UnlockOSThread()
}
