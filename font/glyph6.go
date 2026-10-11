package font

func (t CmapTable) glyph6(r rune) int {
	if len(t.Data) < 10 {
		return 0
	}
	first, count := int(be16(t.Data, 6)), int(be16(t.Data, 8))
	i := int(r) - first
	if i < 0 || i >= count || 10+i*2+2 > len(t.Data) {
		return 0
	}
	return int(be16(t.Data, 10+i*2))
}
