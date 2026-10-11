//go:build darwin

package macos

/*
#cgo CFLAGS: -x objective-c -Wno-deprecated-declarations
#cgo LDFLAGS: -framework Cocoa -framework CoreGraphics -framework CoreVideo
#include <stdlib.h>
#include "backend_darwin.h"
*/
import "C"
import (
	"errors"
	"runtime"
	"unsafe"

	"github.com/gabrielluizsf/antui/backend"
)

// Open creates the window that Options describes.
//
// AppKit insists on the main OS thread. Locking pins this goroutine to the
// thread it is already on, and the CGO layer uses Grand Central Dispatch
// to safely route the NSWindow instantiation to the Apple main thread.
func Open(win backend.Face, opts backend.Options) (*Driver, error) {
	if win == nil {
		return nil, errors.New("antui: macos backend with no window to drive")
	}
	runtime.LockOSThread()

	d := &Driver{win: win}
	cTitle := C.CString(opts.Title)
	defer C.free(unsafe.Pointer(cTitle))

	d.handle = C.antui_d_open(cTitle, C.int(opts.Width), C.int(opts.Height))
	if d.handle == nil {
		runtime.UnlockOSThread()
		return nil, errors.New("antui: could not create the window; on macOS this has to happen on the main thread")
	}
	return d, nil
}
