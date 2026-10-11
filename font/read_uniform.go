package font

// readUniform reads index format 2: every glyph the same size, all with the
// metrics held here.
func (s *strikeReader) readUniform() bool {
	if s.sub+20 > len(s.cblc) {
		return false
	}
	size := int(be32(s.cblc, s.sub+8))
	s.constant, s.hasConst = hMetrics(s.cblc[s.sub+12:]), true
	for g := s.first; g <= s.last; g++ {
		off := s.base + (g-s.first)*size
		s.add(g, off, off+size)
	}
	return true
}
