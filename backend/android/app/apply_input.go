//go:build android

package app

/*
#include <android/configuration.h>
*/
import "C"

// applyInput reads whether there is a keyboard and a touchscreen.
func applyInput(cfg *C.AConfiguration, out *Config) {
	k := C.AConfiguration_getKeyboard(cfg)
	out.HasKeyboard = k != C.ACONFIGURATION_KEYBOARD_ANY && k != C.ACONFIGURATION_KEYBOARD_NOKEYS
	t := C.AConfiguration_getTouchscreen(cfg)
	out.HasTouchscreen = t != C.ACONFIGURATION_TOUCHSCREEN_ANY && t != C.ACONFIGURATION_TOUCHSCREEN_NOTOUCH
}
