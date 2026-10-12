//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/input.h>
#include <android/keycodes.h>
*/
import "C"

// KeyScanCode is the hardware code behind the key, which differs by
// device and is almost never what a program wants.
func (e *InputEvent) KeyScanCode() int { return int(C.AKeyEvent_getScanCode(e.ptr)) }
