//go:build android

package app

// LockOrientation freezes the screen whichever way up it is now.
func LockOrientation(on bool) error {
	if on {
		return requestOrientation(orientationLocked)
	}
	return requestOrientation(orientationUnspecified)
}
