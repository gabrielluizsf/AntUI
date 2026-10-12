//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/input.h>
#include <android/keycodes.h>
*/
import "C"

// The mouse buttons, as a set of bits.
const (
	ButtonPrimary   int32 = C.AMOTION_EVENT_BUTTON_PRIMARY
	ButtonSecondary int32 = C.AMOTION_EVENT_BUTTON_SECONDARY
	ButtonTertiary  int32 = C.AMOTION_EVENT_BUTTON_TERTIARY
	ButtonBack      int32 = C.AMOTION_EVENT_BUTTON_BACK
	ButtonForward   int32 = C.AMOTION_EVENT_BUTTON_FORWARD
)
