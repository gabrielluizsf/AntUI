//go:build android

package camera

import "github.com/gabrielluizsf/antui/backend/android/display"

// uprightRotation works out how far a frame has to be turned.
//
// Two rotations are involved and they go opposite ways. The sensor is
// mounted turned by Orientation; the screen is turned by whatever the user
// has done with the device. A front camera is also mirrored, which changes
// the sign. This is the formula the platform's own documentation gives, and
// it is worth copying rather than deriving.
func uprightRotation(info Info) int {
	screen := 0
	if r, err := display.Rotation(); err == nil {
		screen = r * 90
	}
	if info.Facing == Front {
		return (info.Orientation + screen) % 360
	}
	return (info.Orientation - screen + 360) % 360
}
