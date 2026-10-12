//go:build android

package app

/*
#include <android/configuration.h>
*/
import "C"

// applyDensity reads the screen density, which the platform can also report
// as one of two "unknown" values.
func applyDensity(cfg *C.AConfiguration, out *Config) {
	// Both of those are large numbers rather than zero, so they would be
	// believed.
	d := int(C.AConfiguration_getDensity(cfg))
	if d != C.ACONFIGURATION_DENSITY_ANY && d != C.ACONFIGURATION_DENSITY_NONE && d > 0 {
		out.Density = d
	}
}
