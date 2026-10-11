package canvas

// PaletteIndex is the nearest entry to a colour, or -1 when the canvas has
// no palette.
func (cv *Canvas) PaletteIndex(c Color) int {
	if cv.palette == nil {
		return -1
	}
	if cv.nearest != nil {
		return int(cv.nearest[rgb555(c)])
	}
	best, bestDist := 0, -1
	for i, p := range cv.palette {
		dr := int(c.R()) - int(p.R())
		dg := int(c.G()) - int(p.G())
		db := int(c.B()) - int(p.B())
		d := dr*dr + dg*dg + db*db
		if bestDist < 0 || d < bestDist {
			bestDist, best = d, i
		}
	}
	return best
}
