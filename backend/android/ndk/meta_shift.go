//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/input.h>
#include <android/keycodes.h>
*/
import "C"

// The modifiers this library reads.
const (
	MetaShift Meta = C.AMETA_SHIFT_ON
	MetaAlt   Meta = C.AMETA_ALT_ON
	MetaCtrl  Meta = C.AMETA_CTRL_ON
	MetaSuper Meta = C.AMETA_META_ON
	MetaCaps  Meta = C.AMETA_CAPS_LOCK_ON
)
