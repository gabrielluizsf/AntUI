package font

func (t CmapTable) glyph0(r rune) int {
	if r < 0 || r > 255 || len(t.Data) < 6+256 {
		return 0
	}
	return int(t.Data[6+int(r)])
}
