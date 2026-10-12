//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/input.h>
#include <android/keycodes.h>
*/
import "C"

// The tools a pointer can be.
const (
	ToolUnknown = int(C.AMOTION_EVENT_TOOL_TYPE_UNKNOWN)
	ToolFinger  = int(C.AMOTION_EVENT_TOOL_TYPE_FINGER)
	ToolStylus  = int(C.AMOTION_EVENT_TOOL_TYPE_STYLUS)
	ToolMouse   = int(C.AMOTION_EVENT_TOOL_TYPE_MOUSE)
	ToolEraser  = int(C.AMOTION_EVENT_TOOL_TYPE_ERASER)
)
