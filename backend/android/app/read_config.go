//go:build android

package app

/*
#cgo LDFLAGS: -landroid

#include <android/configuration.h>
#include <android/native_activity.h>
*/
import "C"

// readConfig asks the platform for everything at once.
//
// The configuration has to be read through the asset manager, which is the
// only handle a native activity is given to it — there is no call that
// simply returns it.
func readConfig(act *C.ANativeActivity) Config {
	if act == nil || act.assetManager == nil {
		return Config{}
	}
	cfg := C.AConfiguration_new()
	if cfg == nil {
		return Config{}
	}
	defer C.AConfiguration_delete(cfg)
	C.AConfiguration_fromAssetManager(cfg, act.assetManager)

	out := Config{
		WidthDp:         int(C.AConfiguration_getScreenWidthDp(cfg)),
		HeightDp:        int(C.AConfiguration_getScreenHeightDp(cfg)),
		SmallestWidthDp: int(C.AConfiguration_getSmallestScreenWidthDp(cfg)),
		SDK:             int(C.AConfiguration_getSdkVersion(cfg)),
	}

	// The platform has two ways of saying it does not know the density, and
	// both are large numbers rather than zero, so they would be believed.
	d := int(C.AConfiguration_getDensity(cfg))
	if d != C.ACONFIGURATION_DENSITY_ANY && d != C.ACONFIGURATION_DENSITY_NONE && d > 0 {
		out.Density = d
	}

	switch C.AConfiguration_getOrientation(cfg) {
	case C.ACONFIGURATION_ORIENTATION_PORT:
		out.Orientation = Portrait
	case C.ACONFIGURATION_ORIENTATION_LAND:
		out.Orientation = Landscape
	case C.ACONFIGURATION_ORIENTATION_SQUARE:
		out.Orientation = Square
	}

	out.Night = C.AConfiguration_getUiModeNight(cfg) == C.ACONFIGURATION_UI_MODE_NIGHT_YES
	out.RightToLeft = C.AConfiguration_getLayoutDirection(cfg) == C.ACONFIGURATION_LAYOUTDIR_RTL
	// AConfiguration_getScreenRound would say whether this is a watch face,
	// and is only declared from API 30. Reading it would mean loading the
	// symbol by hand for a device class this library does not target, so it
	// is left out rather than half-done.

	k := C.AConfiguration_getKeyboard(cfg)
	out.HasKeyboard = k != C.ACONFIGURATION_KEYBOARD_ANY && k != C.ACONFIGURATION_KEYBOARD_NOKEYS
	t := C.AConfiguration_getTouchscreen(cfg)
	out.HasTouchscreen = t != C.ACONFIGURATION_TOUCHSCREEN_ANY && t != C.ACONFIGURATION_TOUCHSCREEN_NOTOUCH

	// The language and country are written into two chars each and are not
	// terminated, so they cannot be read as C strings. An unset one comes
	// back as two zero bytes.
	var lang, country [2]C.char
	C.AConfiguration_getLanguage(cfg, &lang[0])
	C.AConfiguration_getCountry(cfg, &country[0])
	out.Language = twoLetters(lang)
	out.Country = twoLetters(country)
	return out
}
