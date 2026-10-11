package svg

// eof reports whether everything has been read.
func (s *pathScanner) eof() bool {
	s.skipSpace()
	return s.pos >= len(s.src)
}
