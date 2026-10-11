//go:build darwin

package macos

/*
#cgo CFLAGS: -x objective-c -Wno-deprecated-declarations
#cgo LDFLAGS: -framework Cocoa -framework CoreGraphics -framework CoreVideo
#include <stdlib.h>
#include "backend_darwin.h"
*/
import "C"
import "unsafe"

// SetClipboard puts text on the system clipboard, reporting whether the
// system took it.
func (d *Driver) SetClipboard(text string) bool {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	C.antui_d_set_clipboard(cText)
	return true
}
