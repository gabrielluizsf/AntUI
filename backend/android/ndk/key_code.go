//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/input.h>
#include <android/keycodes.h>
*/
import "C"

// KeyCode is which key it was, as one of the AKEYCODE constants.
func (e *InputEvent) KeyCode() int { return int(C.AKeyEvent_getKeyCode(e.ptr)) }
