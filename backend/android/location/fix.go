//go:build android

package location

import (
	"time"
)

// Fix is where the device was, and how well it knew.
type Fix struct {
	Latitude, Longitude float64
	// Accuracy is the radius in metres inside which the device believes it
	// is, with about 68 per cent confidence. It is not a promise and it is
	// not a maximum.
	Accuracy float32
	// Altitude is metres above the WGS84 ellipsoid, which is not sea level
	// and can differ from it by a hundred metres. Zero when unknown.
	Altitude float64
	// Speed is metres a second and Bearing is degrees clockwise from north.
	// Both are zero on a device standing still and on one that cannot say.
	Speed   float32
	Bearing float32
	// At is when the position was taken.
	At time.Time
	// Provider is which of the two above answered.
	Provider string
}
