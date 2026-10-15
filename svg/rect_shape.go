package svg

import "github.com/gabrielluizsf/antui/canvas"

// rectShape is a rectangle, which may have its corners rounded. The corners are
// rounded by the arcs SVG describes them with, and a radius bigger than half
// the shorter side is pulled back to it, so a shape never turns inside out.
func rectShape(e *element, warn func(string, ...any)) *canvas.Path {
	x := readLength(warn, e, "x", "x")
	y := readLength(warn, e, "y", "y")
	w, okW := parseLength(e.attr("width"))
	h, okH := parseLength(e.attr("height"))
	if !okW || !okH || w <= 0 || h <= 0 {
		if e.hasAttr("width") && !okW {
			warn("the width %q is not a length, so there is no rectangle to draw", e.attr("width"))
		}
		if e.hasAttr("height") && !okH {
			warn("the height %q is not a length, so there is no rectangle to draw", e.attr("height"))
		}
		// A rectangle with no area is not a shape, which is what SVG says of
		// one with a zero or negative side.
		return nil
	}
	rx, hasRx := readLength(warn, e, "rx", "rx"), e.hasAttr("rx")
	ry, hasRy := readLength(warn, e, "ry", "ry"), e.hasAttr("ry")
	if !hasRx && !hasRy {
		p := canvas.NewPath()
		p.AddRect(canvas.Point{X:x,Y:y}, canvas.Point{X:w,Y:h})
		return p
	}
	// One radius given is the radius on both corners; each is capped at half
	// the side it curves along.
	if !hasRx {
		rx = ry
	}
	if !hasRy {
		ry = rx
	}
	rx = min(rx, w/2)
	ry = min(ry, h/2)
	p := canvas.NewPath()
	// AddRoundRect pulls a radius that is bigger than half its side back to half,
	// which is the rule SVG has, and leaves a radius of no size as a plain
	// rectangle rather than as a shape with no corners at all.
	p.AddRoundRect(canvas.Point{X:x,Y:y}, canvas.Point{X:w,Y:h}, canvas.Point{X:rx,Y:ry})
	return p
}
