package font

// readFlags reads the per-point flags, expanding the repeat counts: a flag
// with bit 3 set is followed by a count of how many more points share it.
func (g *simpleGlyph) readFlags(p int) ([]byte, int, bool) {
	flags := make([]byte, 0, g.points)
	for len(flags) < g.points {
		if p >= len(g.data) {
			return nil, p, false
		}
		flag := g.data[p]
		p++
		flags = append(flags, flag)
		if flag&0x08 != 0 {
			if p >= len(g.data) {
				return nil, p, false
			}
			repeat := int(g.data[p])
			p++
			for range repeat {
				if len(flags) >= g.points {
					break
				}
				flags = append(flags, flag)
			}
		}
	}
	return flags, p, true
}
