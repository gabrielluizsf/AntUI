//go:build android

package speech

// SetRate is how fast: 1 is normal, 2 is twice as fast, 0.5 is half.
func (s *Speaker) SetRate(rate float32) error { return s.tune("setSpeechRate", rate) }
