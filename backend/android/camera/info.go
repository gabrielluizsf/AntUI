//go:build android

package camera

// Info describes one camera.
type Info struct {
	// ID is what [Open] takes. It is an opaque string, not an index.
	ID     string
	Facing Facing
	// Orientation is how far the sensor's image is turned from the device's
	// natural orientation, in degrees clockwise. It is 90 on almost every
	// phone, because the sensor is mounted on its side — which is why a
	// photograph taken by an app that ignores this comes out sideways.
	Orientation int
}
