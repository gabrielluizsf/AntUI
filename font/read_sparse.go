package font

// readSparse reads index format 4: sparse glyphs, each named with its offset,
// and one record past the last of them to measure it.
func (s *strikeReader) readSparse() bool {
	if s.sub+12 > len(s.cblc) {
		return false
	}
	count := int(be32(s.cblc, s.sub+8))
	if count <= 0 || s.sub+12+(count+1)*4 > len(s.cblc) {
		return false
	}
	for i := 0; i < count; i++ {
		at := s.sub + 12 + i*4
		// The record after this one is where this glyph's data ends: its
		// offset, not the id beside it.
		next := s.sub + 12 + (i+1)*4
		s.add(int(be16(s.cblc, at)),
			s.base+int(be16(s.cblc, at+2)),
			s.base+int(be16(s.cblc, next+2)))
	}
	return true
}
