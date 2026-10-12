//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/input.h>
#include <android/keycodes.h>
*/
import "C"

import (
	"unsafe"
)

// WrapQueue takes the pointer a native-activity callback was given.
func WrapQueue(p unsafe.Pointer) *InputQueue {
	if p == nil {
		return nil
	}
	return &InputQueue{ptr: (*C.AInputQueue)(p)}
}
