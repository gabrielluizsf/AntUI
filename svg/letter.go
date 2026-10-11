package svg

// letter reads one command letter.
func (s *pathScanner) letter() byte {
	c := s.src[s.pos]
	s.pos++
	return c
}
