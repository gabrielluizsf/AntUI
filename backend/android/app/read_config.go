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

	applyDensity(cfg, &out)
	applyOrientation(cfg, &out)

	out.Night = C.AConfiguration_getUiModeNight(cfg) == C.ACONFIGURATION_UI_MODE_NIGHT_YES
	out.RightToLeft = C.AConfiguration_getLayoutDirection(cfg) == C.ACONFIGURATION_LAYOUTDIR_RTL
	// AConfiguration_getScreenRound would say whether this is a watch face,
	// and is only declared from API 30. Reading it would mean loading the
	// symbol by hand for a device class this library does not target, so it
	// is left out rather than half-done.

	applyInput(cfg, &out)
	applyLocale(cfg, &out)
	return out
}
