//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/input.h>
#include <android/keycodes.h>
*/
import "C"

// InputEvent is one key press or one touch report. It is only valid until it
// is finished.
type InputEvent struct {
	ptr *C.AInputEvent
}
