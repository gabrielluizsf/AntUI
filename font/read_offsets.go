package font

// readOffsets reads index formats 1 and 3: one offset per glyph in the range,
// and one past the end of it to measure the last. The offsets are 32-bit in
// format 1 and 16-bit in format 3, both added to the base in the header.
func (s *strikeReader) readOffsets(indexFormat int) bool {
	width := 4
	if indexFormat == 3 {
		width = 2
	}
	num := s.last - s.first + 2
	if s.sub+8+num*width > len(s.cblc) {
		return false
	}
	for i := 0; i < num-1; i++ {
		s.add(s.first+i,
			s.base+int(readWidth(s.cblc, s.sub+8+i*width, width)),
			s.base+int(readWidth(s.cblc, s.sub+8+(i+1)*width, width)))
	}
	return true
}
