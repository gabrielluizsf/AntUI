//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid
#include <android/native_window.h>
*/
import "C"

import "unsafe"

// Wrap takes the pointer a native-activity callback was given. It does not
// acquire: the caller owns whatever lifetime the system promised.
func Wrap(p unsafe.Pointer) *Window {
	if p == nil {
		return nil
	}
	return &Window{ptr: (*C.ANativeWindow)(p)}
}
