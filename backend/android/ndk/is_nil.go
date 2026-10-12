//go:build android

package ndk

// IsNil reports whether the device has no such sensor.
func (s Sensor) IsNil() bool { return s.ptr == nil }
