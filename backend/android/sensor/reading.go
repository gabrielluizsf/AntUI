//go:build android

package sensor

import "time"

// Reading is one measurement.
type Reading struct {
	Kind Kind
	// At is when the sensor took it, not when the app read it — on a clock
	// that runs only while the device is awake. The difference between the
	// two is the whole point of batching.
	At time.Duration
	// V is what was measured, kept the full width the platform reports.
	// Three numbers for anything with a direction, one for anything with a
	// quantity, four or five for a rotation.
	V [16]float32
}
