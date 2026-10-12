package svg

// readPatternTile reads where the tile is and how big: nothing at all of the
// drawing, or the whole of it in the drawing's own coordinates. A percentage is
// of the box in the first case and of the drawing in the second, which is the
// same two worlds a gradient's numbers live in.
func (img *Image) readPatternTile(p *pattern, e *element, warn func(string, ...any)) {
	span := [2]float64{img.ViewBox[2], img.ViewBox[3]}
	p.x, p.y = 0, 0
	p.w, p.h = 1, 1
	if !p.boxUnits {
		p.w, p.h = span[0], span[1]
	}
	for _, f := range []struct {
		name string
		dst  *float64
		axis int
	}{
		{"x", &p.x, 0},
		{"y", &p.y, 1},
		{"width", &p.w, 0},
		{"height", &p.h, 1},
	} {
		if !e.hasAttr(f.name) {
			continue
		}
		raw := e.attr(f.name)
		unit := 1.0
		if !p.boxUnits {
			unit = span[f.axis]
		}
		v, ok := gradientLength(raw, unit)
		if !ok {
			warn("the pattern %s %q is not a length, so the tile keeps %g", f.name, raw, *f.dst)
			continue
		}
		*f.dst = v
	}
}
