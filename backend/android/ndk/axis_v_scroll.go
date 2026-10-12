//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/input.h>
#include <android/keycodes.h>
*/
import "C"

// The axes worth naming: the two scroll wheels, and the sticks, hat and
// triggers of a game controller.
const (
	AxisVScroll  = int(C.AMOTION_EVENT_AXIS_VSCROLL)
	AxisHScroll  = int(C.AMOTION_EVENT_AXIS_HSCROLL)
	AxisX        = int(C.AMOTION_EVENT_AXIS_X)
	AxisY        = int(C.AMOTION_EVENT_AXIS_Y)
	AxisZ        = int(C.AMOTION_EVENT_AXIS_Z)
	AxisRZ       = int(C.AMOTION_EVENT_AXIS_RZ)
	AxisHatX     = int(C.AMOTION_EVENT_AXIS_HAT_X)
	AxisHatY     = int(C.AMOTION_EVENT_AXIS_HAT_Y)
	AxisLTrigger = int(C.AMOTION_EVENT_AXIS_LTRIGGER)
	AxisRTrigger = int(C.AMOTION_EVENT_AXIS_RTRIGGER)
	AxisGas      = int(C.AMOTION_EVENT_AXIS_GAS)
	AxisBrake    = int(C.AMOTION_EVENT_AXIS_BRAKE)
)
