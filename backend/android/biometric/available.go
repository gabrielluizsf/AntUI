//go:build android

package biometric

const (
	// Available means the prompt will work.
	Available Availability = iota
	// NoHardware means this device has no sensor.
	NoHardware
	// NoneEnrolled means it has one and the user has not registered a
	// fingerprint or a face on it.
	NoneEnrolled
	// Unavailable means the hardware is there and busy or broken.
	Unavailable
	// TooOld means the platform predates the prompt: below API 28 there is
	// no BiometricPrompt, only the fingerprint API it replaced.
	TooOld
)
