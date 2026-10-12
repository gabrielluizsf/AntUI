//go:build android

package sensor

var kindNames = map[Kind]string{
	Accelerometer: "accelerometer", MagneticField: "magnetic field",
	Gyroscope: "gyroscope", Light: "light", Pressure: "pressure",
	Proximity: "proximity", Humidity: "humidity", Temperature: "temperature",
	Gravity: "gravity", LinearAcceleration: "linear acceleration",
	RotationVector: "rotation vector", GameRotationVector: "game rotation vector",
	GeomagneticRotation: "geomagnetic rotation", StepDetector: "step detector",
	StepCounter: "step counter", SignificantMotion: "significant motion",
}
