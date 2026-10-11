//go:build darwin

package macos

/*
#cgo CFLAGS: -x objective-c -Wno-deprecated-declarations
#cgo LDFLAGS: -framework Cocoa -framework CoreGraphics -framework CoreVideo
#include <stdlib.h>
#include "backend_darwin.h"
*/
import "C"

// SystemDark is the appearance macOS paints its own interface in, read from
// the defaults the system keeps. macOS always has one to give — a Mac older
// than dark mode is a light one — so the answer is always known.
func (d *Driver) SystemDark() (bool, bool) { return C.antui_d_system_dark() == 1, true }
