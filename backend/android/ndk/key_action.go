//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/input.h>
#include <android/keycodes.h>
*/
import "C"

// KeyAction is whether the key went down, came up, or is repeating.
func (e *InputEvent) KeyAction() int { return int(C.AKeyEvent_getAction(e.ptr)) }
