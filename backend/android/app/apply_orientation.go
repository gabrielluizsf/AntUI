//go:build android

package app

/*
#include <android/configuration.h>
*/
import "C"

// applyOrientation reads which way up the screen is.
func applyOrientation(cfg *C.AConfiguration, out *Config) {
	switch C.AConfiguration_getOrientation(cfg) {
	case C.ACONFIGURATION_ORIENTATION_PORT:
		out.Orientation = Portrait
	case C.ACONFIGURATION_ORIENTATION_LAND:
		out.Orientation = Landscape
	case C.ACONFIGURATION_ORIENTATION_SQUARE:
		out.Orientation = Square
	}
}
