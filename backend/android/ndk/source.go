//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/input.h>
#include <android/keycodes.h>
*/
import "C"

// Source is what produced it: a touchscreen, a mouse, a joystick. The value
// is a set of bits, so it is tested with [Source.Is] rather than compared.
func (e *InputEvent) Source() Source { return Source(C.AInputEvent_getSource(e.ptr)) }
