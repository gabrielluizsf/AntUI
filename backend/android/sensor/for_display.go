//go:build android

package sensor

// ForDisplay turns a reading from the device's frame into the one the app is
// drawing in.
//
// This is the thing everyone gets wrong once. A sensor reports in a frame
// fixed to the *hardware*: X across the device as it was built. The screen
// rotates and the sensor does not, so on a phone held sideways the
// accelerometer's X runs down the screen and a ball rolls the wrong way.
// The rotation to pass is [antui/backend/android/display.Rotation].
func (r Reading) ForDisplay(rotation int) Reading {
	r.V[0], r.V[1] = remap(r.V[0], r.V[1], rotation)
	return r
}
