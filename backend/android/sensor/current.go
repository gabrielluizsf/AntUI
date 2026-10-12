//go:build android

package sensor

// Current is the most recent reading, whether or not it is new.
func (s *Stream) Current(k Kind) (Reading, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.latest[k]
	return r, ok
}
