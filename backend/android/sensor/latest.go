//go:build android

package sensor

// Latest is the most recent reading from one sensor, and whether there has
// been a new one since the last call. It is what a frame loop wants: the
// current value, with no queue to drain.
func (s *Stream) Latest(k Kind) (Reading, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.latest[k]
	fresh := s.fresh[k]
	s.fresh[k] = false
	return r, ok && fresh
}
