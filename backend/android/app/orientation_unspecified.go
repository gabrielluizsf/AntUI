//go:build android

package app

// The values android.content.pm.ActivityInfo uses for a requested
// orientation. Only the four worth asking for are here.
const (
	orientationUnspecified = -1
	orientationLandscape   = 0
	orientationPortrait    = 1
	orientationLocked      = 14 // whichever way it is now, and stay there
)
