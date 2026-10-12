package svg

import "strings"

// readPattern reads where one `<pattern>`'s tile is and how it is moved —
// everything about it except the picture it holds, which is built afterwards by
// [Image.readPatterns]. What it does not say is the SVG default, which is the
// box of the shape it paints for the tile and the drawing's own coordinates for
// what is inside it.
func (img *Image) readPattern(id string, e *element) {
	warn := func(format string, args ...any) { img.warn(e, format, args...) }
	p := img.patterns[id]
	p.boxUnits = !strings.EqualFold(strings.TrimSpace(e.attr("patternUnits")), "userspaceonuse")
	p.contentBox = strings.EqualFold(strings.TrimSpace(e.attr("patternContentUnits")), "objectboundingbox")
	img.readPatternTile(p, e, warn)
	img.readPatternTransform(p, e, warn)
	img.readPatternHref(p, e, warn)
	// Read after any href above, so that what this element says about its own
	// room is what stands rather than what the pattern it borrowed from said.
	p.visible = overflowIn(e)
	img.readPatternViewBox(p, e, warn)
	if p.w <= 0 || p.h <= 0 {
		warn("the pattern tile is %g by %g, which tiles nothing, so the shape is painted without it", p.w, p.h)
	}
}

// readPatternTransform reads `patternTransform`, the transform that moves the
// tile rather than the shape it paints.
func (img *Image) readPatternTransform(p *pattern, e *element, warn func(string, ...any)) {
	v := strings.TrimSpace(e.attr("patternTransform"))
	if v == "" {
		return
	}
	if m, ok := parseTransform(v); ok {
		p.transform = m
	} else {
		warn("the pattern transform %q is not one this package can read, so the tile stays where it is", v)
	}
}
