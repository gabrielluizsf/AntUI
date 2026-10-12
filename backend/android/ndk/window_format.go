//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid
#include <android/native_window.h>
*/
import "C"

// Format is the layout of the pixels in the buffer.
func (w *Window) Format() Format { return Format(C.ANativeWindow_getFormat(w.ptr)) }
