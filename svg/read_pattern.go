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

	// The tile before anything says where it is: nothing at all of the drawing,
	// or the whole of it in the drawing's own coordinates. A percentage is of
	// the box in the first case and of the drawing in the second, which is the
	// same two worlds a gradient's numbers live in.
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
	if v := strings.TrimSpace(e.attr("patternTransform")); v != "" {
		if m, ok := parseTransform(v); ok {
			p.transform = m
		} else {
			warn("the pattern transform %q is not one this package can read, so the tile stays where it is", v)
		}
	}
	for _, name := range []string{"href", "xlink:href"} {
		if ref, ok := hashRef(e.attr(name)); ok && ref != "" {
			// Follow the reference to another pattern, the same way gradients
			// follow href references. The referenced pattern must be named with an
			// id and be within this drawing.
			id := strings.TrimPrefix(ref, "#")
			if refPattern, ok := img.patterns[id]; ok {
				// Use the referenced pattern's content and transform,
				// inheriting its boxUnits and contentBox settings.
				p.content = refPattern.content
				p.transform = refPattern.transform
				p.boxUnits = refPattern.boxUnits
				p.contentBox = refPattern.contentBox
				warn("the pattern takes its content from #%s", id)
			} else {
				warn("the pattern references #%s, which this drawing does not have, so its content is empty", id)
			}
			break
		}
	}
	// Read after any href above, so that what this element says about its own
	// room is what stands rather than what the pattern it borrowed from said.
	p.visible = overflowIn(e)
	if e.hasAttr("viewBox") {
		vbStr := strings.TrimSpace(e.attr("viewBox"))
		if vb, ok := parseViewBox(vbStr); ok {
			// viewBox establishes the pattern's coordinate system, overriding
			// x, y, width, height and patternUnits/patternContentUnits.
			p.x, p.y, p.w, p.h = vb[0], vb[1], vb[2], vb[3]
			p.boxUnits = true
			// patternContentUnits is also ignored per spec when viewBox is given,
			// but we keep the current behavior for backward compatibility.
			if e.hasAttr("patternContentUnits") {
				p.contentBox = strings.EqualFold(
					strings.TrimSpace(e.attr("patternContentUnits")), "objectboundingbox")
			}
			warn("the pattern has a viewBox, which sets the tile coordinates; " +
				"x, y, width and height attributes are ignored")
		} else {
			warn("the pattern has an invalid viewBox %q, so its picture is drawn where it was written", vbStr)
		}
	} else if e.hasAttr("viewBox") {
		// viewBox present but invalid — fall through to warning below
	}
	if p.w <= 0 || p.h <= 0 {
		warn("the pattern tile is %g by %g, which tiles nothing, so the shape is painted without it", p.w, p.h)
	}
}
