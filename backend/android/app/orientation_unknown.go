//go:build android

package app

// The orientations. Square is a shape no phone has had in a decade and is
// here because the platform still reports it.
const (
	OrientationUnknown Orientation = iota
	Portrait
	Landscape
	Square
)
