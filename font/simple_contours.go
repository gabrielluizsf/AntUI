package font

// simpleContours reads a glyph that is its own outline.
func simpleContours(data []byte, contours int) []Contour {
	if contours == 0 || 10+contours*2+2 > len(data) {
		return nil
	}
	g := &simpleGlyph{data: data, ends: make([]int, contours)}
	for i := range contours {
		g.ends[i] = int(be16(data, 10+i*2))
	}
	g.points = g.ends[contours-1] + 1
	if g.points <= 0 || g.points > 100000 {
		return nil
	}

	p := 10 + contours*2
	instructions := int(be16(data, p))
	p += 2 + instructions
	if p > len(data) {
		return nil
	}

	var ok bool
	g.flags, p, ok = g.readFlags(p)
	if !ok {
		return nil
	}
	xs, p, ok := g.readCoords(p, 0x02, 0x10)
	if !ok {
		return nil
	}
	ys, _, ok := g.readCoords(p, 0x04, 0x20)
	if !ok {
		return nil
	}
	return g.buildContours(xs, ys)
}
