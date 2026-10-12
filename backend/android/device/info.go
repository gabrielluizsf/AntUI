//go:build android

package device

// Info is the device, as far as an app is allowed to know it.
type Info struct {
	// Manufacturer, Brand and Model are what is printed on the box and in
	// the settings screen: "Google", "google", "Pixel 8".
	Manufacturer string
	Brand        string
	Model        string
	// Device and Product are the internal names, which is what a bug report
	// and a compatibility list are keyed on.
	Device  string
	Product string
	// SDK is the API level, and Release is the version a person would say:
	// 34 and "14".
	SDK     int
	Release string
	// Locale is the language and region the user chose, as "pt-BR".
	Locale string
	// TimeZone is the zone's id, as "America/Sao_Paulo".
	TimeZone string
	// Emulator is a guess, from the build's own fingerprint. It is right far
	// more often than not and should never be the basis of anything but a
	// diagnostic.
	Emulator bool
}
