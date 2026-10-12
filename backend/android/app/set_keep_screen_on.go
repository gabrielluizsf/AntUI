//go:build android

package app

// SetKeepScreenOn stops the display dimming and locking while the app is in
// front. It is a window flag rather than a Java call, so it needs neither
// the UI thread nor the bridge.
func SetKeepScreenOn(on bool) {
	if on {
		SetWindowFlags(FlagKeepScreenOn, 0)
	} else {
		SetWindowFlags(0, FlagKeepScreenOn)
	}
}
