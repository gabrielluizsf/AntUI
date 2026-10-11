//go:build darwin

package macos

/*
#cgo CFLAGS: -x objective-c -Wno-deprecated-declarations
#cgo LDFLAGS: -framework Cocoa -framework CoreGraphics -framework CoreVideo
#include <stdlib.h>
#include "backend_darwin.h"
*/
import "C"

// droppedFiles reads the paths of the drop that has just been reported. The
// window holds them until the next drop, so this is safe to do while the
// event is being decoded.
func (d *Driver) droppedFiles(handle *C.antui_d_window, count int) []string {
	if handle == nil || count <= 0 {
		return nil
	}
	held := int(C.antui_d_drop_count(handle))
	if held < count {
		count = held
	}
	files := make([]string, 0, count)
	for i := range count {
		if path := C.antui_d_drop_path(handle, C.int(i)); path != nil {
			files = append(files, C.GoString(path))
		}
	}
	return files
}
