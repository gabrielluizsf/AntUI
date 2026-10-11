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

// Clipboard is what the system clipboard holds: its text, and the files it
// names.
func (d *Driver) Clipboard() (text string, files []string) {
	if held := C.antui_d_clipboard_text(); held != nil {
		text = C.GoString(held)
		C.free(unsafe.Pointer(held))
	}

	for i := range int(C.antui_d_clipboard_count()) {
		if path := C.antui_d_clipboard_path(C.int(i)); path != nil {
			files = append(files, C.GoString(path))
		}
	}
	return text, files
}
