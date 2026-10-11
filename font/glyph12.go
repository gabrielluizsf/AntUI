package font

func (t CmapTable) glyph12(r rune) int {
	for i := range t.Groups {
		g := 16 + i*12
		if g+12 > len(t.Data) {
			return 0
		}
		start, end := be32(t.Data, g), be32(t.Data, g+4)
		if uint32(r) < start {
			return 0
		}
		if uint32(r) > end {
			continue
		}
		return int(be32(t.Data, g+8) + uint32(r) - start)
	}
	return 0
}
