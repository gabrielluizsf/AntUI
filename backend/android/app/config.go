//go:build android

package app

// Config is the device configuration the app is running under. All of it can
// change while the app is running — a rotation, a language, a theme, a
// device unfolding — and every change arrives as a [ConfigChanged] event
// with this already updated.
type Config struct {
	// Density is dots per inch as Android buckets it: 160 means one pixel
	// per device-independent pixel. Dividing by 160 gives the scale.
	Density int
	// Orientation is which way up the screen is.
	Orientation Orientation
	// Night is whether the device is in dark mode. It is the one thing a
	// theme has to follow, and following it is what makes an app look like
	// it belongs on the device.
	Night bool
	// Language and Country are two letters each, lower and upper case:
	// "pt", "BR". Empty when the platform will not say.
	Language, Country string
	// RightToLeft is the layout direction the user's language wants.
	RightToLeft bool
	// WidthDp and HeightDp are the window in device-independent pixels, and
	// SmallestWidthDp is the smaller of the two whichever way the device is
	// turned — which is the number a layout should choose on, because it
	// does not change when the screen rotates.
	WidthDp, HeightDp int
	SmallestWidthDp   int
	// Keyboard and Touchscreen are what the device has attached.
	HasKeyboard    bool
	HasTouchscreen bool
	// SDK is the platform's API level.
	SDK int
}
