//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/input.h>
#include <android/keycodes.h>
*/
import "C"

// The sources worth naming. Each is a class of device or'd with the kind of
// data it produces, which is why they overlap and are tested as bits.
const (
	SourceKeyboard    Source = C.AINPUT_SOURCE_KEYBOARD
	SourceDpad        Source = C.AINPUT_SOURCE_DPAD
	SourceGamepad     Source = C.AINPUT_SOURCE_GAMEPAD
	SourceTouchscreen Source = C.AINPUT_SOURCE_TOUCHSCREEN
	SourceMouse       Source = C.AINPUT_SOURCE_MOUSE
	SourceStylus      Source = C.AINPUT_SOURCE_STYLUS
	SourceTrackball   Source = C.AINPUT_SOURCE_TRACKBALL
	SourceTouchpad    Source = C.AINPUT_SOURCE_TOUCHPAD
	SourceJoystick    Source = C.AINPUT_SOURCE_JOYSTICK
)
