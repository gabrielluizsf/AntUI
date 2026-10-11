package font

// add places one glyph's data, as long as it lies inside the table.
func (s *strikeReader) add(g, off, end int) {
	if g < s.first || g > s.last || off < 0 || off >= end || end > len(s.cbdt) {
		return
	}
	s.out[g] = colorBitmap{off: off, length: end - off, imageFormat: s.imageFormat}
}
