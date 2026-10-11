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

// Driver is one window on macOS. It owns the platform connections the way
// the window core owns the input state.
type Driver struct {
	handle *C.antui_d_window
	win    backend.Face
	events [64]C.antui_d_event
}
