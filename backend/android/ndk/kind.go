//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/input.h>
#include <android/keycodes.h>
*/
import "C"

// Kind is whether this is a key or a motion.
func (e *InputEvent) Kind() int { return int(C.AInputEvent_getType(e.ptr)) }
