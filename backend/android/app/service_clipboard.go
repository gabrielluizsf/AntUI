//go:build android

package app

// The names of the system services this library reaches for. Android has
// dozens more; these are the ones with a package above them.
const (
	ServiceClipboard    = "clipboard"
	ServiceNotification = "notification"
	ServiceConnectivity = "connectivity"
	ServiceVibrator     = "vibrator"
	ServiceWindow       = "window"
	ServicePower        = "power"
	ServiceInputMethod  = "input_method"
	ServiceDisplay      = "display"
	ServiceSensor       = "sensor"
	ServiceLocation     = "location"
)
