//go:build android

package device

// looksEmulated is a guess and is named like one.
func looksEmulated(fingerprint, model, manufacturer string) bool {
	for _, sign := range []string{"generic", "unknown", "emulator", "sdk_gphone",
		"Android SDK built for", "Genymotion", "vbox"} {
		if contains(fingerprint, sign) || contains(model, sign) || contains(manufacturer, sign) {
			return true
		}
	}
	return false
}
