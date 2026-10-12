//go:build android

package ndk

import (
	"unsafe"
)

// Pointer is the raw handle, for passing between packages that each have
// their own idea of what a C type is.
func (w *Window) Pointer() unsafe.Pointer {
	if w == nil {
		return nil
	}
	return unsafe.Pointer(w.ptr)
}
