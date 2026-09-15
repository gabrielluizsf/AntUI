//go:build android

package app

/*
#cgo LDFLAGS: -landroid

#include <android/configuration.h>
#include <android/native_activity.h>
*/
import "C"

// Orientation is which way up the screen is.
type Orientation int

// The orientations. Square is a shape no phone has had in a decade and is
// here because the platform still reports it.
const (
	OrientationUnknown Orientation = iota
	Portrait
	Landscape
	Square
)

// String names the orientation the way Android does.
func (o Orientation) String() string {
	switch o {
	case Portrait:
		return "portrait"
	case Landscape:
		return "landscape"
	case Square:
		return "square"
	}
	return "unknown"
}

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

// Scale is the factor a size in points is multiplied by, or 0 when the
// platform will not say.
func (c Config) Scale() float64 {
	if c.Density <= 0 {
		return 0
	}
	return float64(c.Density) / 160
}

// Locale is the language and country joined the way a tag is written:
// "pt-BR", or just "pt" when there is no country.
func (c Config) Locale() string {
	switch {
	case c.Language == "":
		return ""
	case c.Country == "":
		return c.Language
	}
	return c.Language + "-" + c.Country
}

// Config is the configuration the app is running under.
func (a *App) Config() Config {
	if a == nil {
		return Config{}
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.config
}

// Density is the screen's density in dots per inch. It is [Config.Density],
// kept as its own call because it is what the window layer asks for on every
// configuration change.
func (a *App) Density() int {
	if a == nil {
		return 0
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.config.Density
}

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

func twoLetters(b [2]C.char) string {
	if b[0] == 0 {
		return ""
	}
	return string([]byte{byte(b[0]), byte(b[1])})
}

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

// SetWindowFlags adds and removes flags on the activity's window. It may be
// called from any thread; the platform posts it to the UI thread itself.
func SetWindowFlags(add, remove uint32) {
	stateMu.Lock()
	act := activity
	stateMu.Unlock()
	if act != nil {
		C.ANativeActivity_setWindowFlags(act, C.uint32_t(add), C.uint32_t(remove))
	}
}

// Flags for the on-screen keyboard, as the platform numbers them.
const (
	// KeyboardImplicit asks for the keyboard the way tapping a text field
	// does: the system may decline, on a device with a real keyboard
	// attached.
	KeyboardImplicit = 0x0001
	// KeyboardForced insists. The user cannot dismiss it with the back
	// button, so it is the wrong choice for almost everything.
	KeyboardForced = 0x0002
	// KeyboardImplicitOnly hides the keyboard only if it was shown
	// implicitly, and KeyboardNotAlways leaves it up if the user asked for
	// it by hand.
	KeyboardImplicitOnly = 0x0001
	KeyboardNotAlways    = 0x0002
)

// ShowKeyboard asks for the on-screen keyboard.
//
// This is the whole of what the NDK offers: it can be asked for and asked to
// go away. **How tall it is, and what is typed on it, are not here** — the
// height is a window inset and the text goes through an input connection,
// and both of those are Java. Until the JNI layer lands, a native app can
// raise the keyboard and will not hear a word of what it produces.
func ShowKeyboard(flags uint32) {
	stateMu.Lock()
	act := activity
	stateMu.Unlock()
	if act != nil {
		C.ANativeActivity_showSoftInput(act, C.uint32_t(flags))
	}
}

// HideKeyboard puts the on-screen keyboard away.
func HideKeyboard(flags uint32) {
	stateMu.Lock()
	act := activity
	stateMu.Unlock()
	if act != nil {
		C.ANativeActivity_hideSoftInput(act, C.uint32_t(flags))
	}
}
