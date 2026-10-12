//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/input.h>
#include <android/keycodes.h>
*/
import "C"

// The motion actions, after masking. DOWN and UP are the first finger going
// down and the last coming up; POINTER_DOWN and POINTER_UP are the ones in
// between, and they carry which finger in the top bits — see
// [InputEvent.MotionIndex].
const (
	MotionDown        = int(C.AMOTION_EVENT_ACTION_DOWN)
	MotionUp          = int(C.AMOTION_EVENT_ACTION_UP)
	MotionMove        = int(C.AMOTION_EVENT_ACTION_MOVE)
	MotionCancel      = int(C.AMOTION_EVENT_ACTION_CANCEL)
	MotionOutside     = int(C.AMOTION_EVENT_ACTION_OUTSIDE)
	MotionPointerDown = int(C.AMOTION_EVENT_ACTION_POINTER_DOWN)
	MotionPointerUp   = int(C.AMOTION_EVENT_ACTION_POINTER_UP)
	MotionHoverMove   = int(C.AMOTION_EVENT_ACTION_HOVER_MOVE)
	MotionScroll      = int(C.AMOTION_EVENT_ACTION_SCROLL)
	MotionHoverEnter  = int(C.AMOTION_EVENT_ACTION_HOVER_ENTER)
	MotionHoverExit   = int(C.AMOTION_EVENT_ACTION_HOVER_EXIT)
	MotionButtonPress = int(C.AMOTION_EVENT_ACTION_BUTTON_PRESS)
	MotionButtonUp    = int(C.AMOTION_EVENT_ACTION_BUTTON_RELEASE)
)
