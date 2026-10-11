package font

// readCoords reads one axis of the point coordinates, which are deltas in one
// of three widths depending on the flags: a byte, a byte that means "the same
// as last time", or a word.
func (g *simpleGlyph) readCoords(p int, shortBit, sameBit byte) ([]float64, int, bool) {
	out := make([]float64, g.points)
	v := 0.0
	for i, flag := range g.flags {
		switch {
		case flag&shortBit != 0:
			if p >= len(g.data) {
				return nil, p, false
			}
			d := float64(g.data[p])
			p++
			if flag&sameBit == 0 {
				d = -d
			}
			v += d
		case flag&sameBit == 0:
			if p+2 > len(g.data) {
				return nil, p, false
			}
			v += float64(int16(be16(g.data, p)))
			p += 2
		}
		out[i] = v
	}
	return out, p, true
}
