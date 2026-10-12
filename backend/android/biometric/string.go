//go:build android

package biometric

// String says why biometrics can or cannot be asked for, in a line that
// can be shown to a person.
func (a Availability) String() string {
	switch a {
	case Available:
		return "available"
	case NoHardware:
		return "no hardware"
	case NoneEnrolled:
		return "nothing enrolled"
	case Unavailable:
		return "unavailable"
	}
	return "too old"
}
