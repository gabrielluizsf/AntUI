//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/input.h>
#include <android/keycodes.h>
*/
import "C"

// KeyRepeat is how many times the key has repeated. It is 0 on the first
// press, which is how a repeat is told from a new one.
func (e *InputEvent) KeyRepeat() int { return int(C.AKeyEvent_getRepeatCount(e.ptr)) }
