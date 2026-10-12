//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/input.h>
#include <android/keycodes.h>
*/
import "C"

// The two kinds of event.
const (
	EventKey    = int(C.AINPUT_EVENT_TYPE_KEY)
	EventMotion = int(C.AINPUT_EVENT_TYPE_MOTION)
)
