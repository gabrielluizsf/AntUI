package font

// readSparseUniform reads index format 5: sparse glyphs with one size and one
// set of metrics for the whole range — imageSize, the metrics, the count, then
// the ids.
func (s *strikeReader) readSparseUniform() bool {
	if s.sub+24 > len(s.cblc) {
		return false
	}
	size := int(be32(s.cblc, s.sub+8))
	count := int(be32(s.cblc, s.sub+20))
	if count <= 0 || s.sub+24+count*2 > len(s.cblc) {
		return false
	}
	s.constant, s.hasConst = hMetrics(s.cblc[s.sub+12:]), true
	for i := range count {
		off := s.base + i*size
		s.add(int(be16(s.cblc, s.sub+24+i*2)), off, off+size)
	}
	return true
}
