//go:build android

package app

// Window flags, as the platform numbers them. Only the ones this library has
// a use for are here.
const (
	// FlagFullscreen hides the status bar. It does not hide the navigation
	// bar: that is immersive mode, which has no NDK call and arrives with
	// the JNI layer.
	FlagFullscreen = 0x00000400
	// FlagKeepScreenOn stops the display dimming and locking while the app
	// is in front.
	FlagKeepScreenOn = 0x00000080
	// FlagLayoutInScreen and FlagLayoutNoLimits put the window behind the
	// system bars rather than inside them.
	FlagLayoutInScreen = 0x00000100
	FlagLayoutNoLimits = 0x00000200
)
