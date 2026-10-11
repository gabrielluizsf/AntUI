package svg

// atLetter reports whether the next thing is a command rather than a number.
func (s *pathScanner) atLetter() bool {
	s.skipSpace()
	if s.pos >= len(s.src) {
		return false
	}
	c := s.src[s.pos]
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}
