//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/input.h>
#include <android/keycodes.h>
*/
import "C"

// The key actions.
const (
	KeyDown = int(C.AKEY_EVENT_ACTION_DOWN)
	KeyUp   = int(C.AKEY_EVENT_ACTION_UP)
)
