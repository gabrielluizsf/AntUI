//go:build android

package speech

// SetPitch is how high: 1 is normal.
func (s *Speaker) SetPitch(pitch float32) error { return s.tune("setPitch", pitch) }
