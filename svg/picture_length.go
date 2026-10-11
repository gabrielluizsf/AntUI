package svg

import "strings"

// pictureLength is one of the four numbers the box is written with: a length
// in the drawing's own units, or a percentage of the drawing — which is the
// viewBox, the same thing a percentage of a gradient is a fraction of. An
// attribute that is not there is zero, and one that is written as something
// this cannot read says so and counts as nothing, the same as every other
// length the drawing writes. axis is which side of the drawing a percentage is
// a fraction of: 0 for across, 1 for down.
func (img *Image) pictureLength(warn func(string, ...any), e *element, attr string, axis int) float64 {
	if !e.hasAttr(attr) {
		return 0
	}
	raw := strings.TrimSpace(e.attr(attr))
	if p, ok := strings.CutSuffix(raw, "%"); ok {
		v, err := parseNumber(p)
		if err != nil {
			warn("the %s %q is not a length, so it is taken as none", attr, e.attr(attr))
			return 0
		}
		span := img.ViewBox[2]
		if axis == 1 {
			span = img.ViewBox[3]
		}
		return v / 100 * span
	}
	v, ok := parseLength(raw)
	if !ok {
		warn("the %s %q is not a length, so it is taken as none", attr, e.attr(attr))
	}
	return v
}
