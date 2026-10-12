//go:build android

package biometric

// Availability is whether the device can ask for a fingerprint or a face,
// and if not, why not. The distinction matters: hardware that is missing is
// permanent, and a user who has enrolled nothing can be sent to the settings
// screen to fix it.
type Availability int
