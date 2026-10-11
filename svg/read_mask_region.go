package svg

import "strings"

// readMaskRegion reads where a `<mask>` reaches: `maskUnits` says which of the
// two coordinates the four numbers are written in, and `x`, `y`, `width` and
// `height` are the rectangle itself. What is not written is the SVG default of
// the box grown by a tenth all round, which is the same numbers either way —
// a fraction of the box, or of the drawing when the region is written in the
// drawing's own coordinates.
func (img *Image) readMaskRegion(e *element, warn func(string, ...any)) maskRegion {
	units := strings.ToLower(strings.TrimSpace(e.attr("maskUnits")))
	r := maskRegion{boxUnits: units != "userspaceonuse", x: -0.1, y: -0.1, w: 1.2, h: 1.2}
	vb := img.ViewBox
	if !r.boxUnits {
		r.x, r.y, r.w, r.h = -0.1*vb[2], -0.1*vb[3], 1.2*vb[2], 1.2*vb[3]
	}
	for _, f := range []struct {
		name string
		dst  *float64
		span float64
	}{
		{"x", &r.x, vb[2]},
		{"y", &r.y, vb[3]},
		{"width", &r.w, vb[2]},
		{"height", &r.h, vb[3]},
	} {
		if !e.hasAttr(f.name) {
			continue
		}
		v, ok := maskRegionValue(e.attr(f.name), r.boxUnits, f.span)
		if !ok {
			warn("the mask region %s %q is not a number this package can read, so it is taken as the whole of it", f.name, e.attr(f.name))
			continue
		}
		*f.dst = v
	}
	return r
}
