package svg

// flag reads the zero or one of the arc command, which SVG allows to be written
// without any separator at all, as in `a1 1 0 011 1`.
func (s *pathScanner) flag() (bool, bool) {
	s.skipSpace()
	if s.pos >= len(s.src) {
		return false, false
	}
	switch s.src[s.pos] {
	case '0':
		s.pos++
		return false, true
	case '1':
		s.pos++
		return true, true
	}
	return false, false
}
