//go:build android

package app

// SetOrientation asks the system to turn the screen and keep it there.
// [OrientationUnknown] gives the choice back to the device.
func SetOrientation(o Orientation) error {
	want := orientationUnspecified
	switch o {
	case Portrait:
		want = orientationPortrait
	case Landscape:
		want = orientationLandscape
	}
	return requestOrientation(want)
}
