//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/input.h>
#include <android/keycodes.h>
*/
import "C"

// Device is the input device's id, which is how two gamepads are told apart.
func (e *InputEvent) Device() int { return int(C.AInputEvent_getDeviceId(e.ptr)) }
