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

// Pump handles whatever happened this frame, folding it into the window
// through win.Push.
func (d *Driver) Pump(win backend.Face) {
	if d.handle == nil {
		return
	}
	// Loop until the queue comes back short, so a burst bigger than the
	// buffer is drained rather than left to arrive a frame late.
	for {
		n := int(C.antui_d_pump(d.handle, &d.events[0], C.int(len(d.events))))
		for i := range n {
			win.Push(d.decode(win, &d.events[i]))
		}
		if n < len(d.events) {
			return
		}
	}
}
