//go:build android

package sensor

// The sensors worth naming. The three groups are worth telling apart.
//
// The first measure something directly and report it raw, noise and all.
// The second are computed by the platform out of the first — a rotation
// vector is the accelerometer, the gyroscope and the compass fused, and is
// steadier than any of them and slower to react. The third report an event
// rather than a value.
const (
	// Measured directly.
	Accelerometer Kind = 1  // metres per second squared, gravity included
	MagneticField Kind = 2  // microtesla
	Gyroscope     Kind = 4  // radians per second
	Light         Kind = 5  // lux
	Pressure      Kind = 6  // hectopascal
	Proximity     Kind = 8  // centimetres, and on most phones only 0 or "far"
	Humidity      Kind = 12 // per cent
	Temperature   Kind = 13 // degrees celsius

	// Computed by the platform.
	Gravity            Kind = 9  // which way is down, without the movement
	LinearAcceleration Kind = 10 // the movement, without gravity
	RotationVector     Kind = 11 // orientation as a quaternion, using the compass
	GameRotationVector Kind = 15 // the same without the compass: no true north,
	// and no lurching when a magnet goes past
	GeomagneticRotation Kind = 20 // the same using the compass and no gyroscope

	// Events rather than values.
	StepDetector      Kind = 18 // one reading per step
	StepCounter       Kind = 19 // steps since the device booted, in V[0]
	SignificantMotion Kind = 17
)
