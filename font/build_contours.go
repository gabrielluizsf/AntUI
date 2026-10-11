package font

// buildContours cuts the read points into one Contour per end.
func (g *simpleGlyph) buildContours(xs, ys []float64) []Contour {
	out := make([]Contour, 0, len(g.ends))
	start := 0
	for _, end := range g.ends {
		if end < start || end >= g.points {
			break
		}
		c := Contour{
			X:  xs[start : end+1],
			Y:  ys[start : end+1],
			On: make([]bool, end+1-start),
		}
		for i := range c.On {
			c.On[i] = g.flags[start+i]&0x01 != 0
		}
		out = append(out, c)
		start = end + 1
	}
	return out
}
