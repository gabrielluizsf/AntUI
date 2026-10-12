//go:build android

package ndk

// SensorEvent is one reading.
type SensorEvent struct {
	// Type is which kind of sensor it came from.
	Type int
	// ID identifies the sensor itself, which matters when two of a kind are
	// enabled.
	ID int
	// At is the reading's own timestamp, in nanoseconds on a clock that runs
	// while the device is awake. It is not the time the app read it, and the
	// difference is the whole point of batching.
	At int64
	// V is what was measured. The platform's array is sixteen long and it is
	// kept whole: three numbers for anything that measures a direction, one
	// for anything that measures a quantity, four or five for a rotation.
	V [16]float32
}
