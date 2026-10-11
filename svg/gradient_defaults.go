package svg

// gradientDefaults is what a gradient is before it says anything, and what it
// falls back to for every attribute it leaves out. The default is a linear
// gradient running left to right across the shape and, in fractions of that shape,
// a circle at the middle of it half as wide as it is.
func gradientDefaults(radial bool) *gradient {
	g := &gradient{radial: radial, spread: spreadPad, set: map[string]bool{}, pct: map[string]pctRead{}}
	if radial {
		g.cx, g.cy, g.r = 0.5, 0.5, 0.5
		return g
	}
	g.x2 = 1
	return g
}
