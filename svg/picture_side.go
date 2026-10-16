package svg

import "strings"

// pictureSide is one of the two sides of an `<image>` box: the length the
// drawing wrote, and whether the side was left open for the picture to say —
// an attribute that is not there or is written as `auto` is the spec's way of
// asking for the picture's own size, or, across from a side that was written,
// for its proportion. A length that cannot be read is zero and said out loud,
// the same as every other length the drawing writes.
func (img *Image) pictureSide(warn func(string, ...any), e *element, attr string, axis int) (float64, bool) {
	if !e.hasAttr(attr) || strings.EqualFold(strings.TrimSpace(e.attr(attr)), "auto") {
		return 0, true
	}
	return img.pictureLength(warn, e, attr, axis), false
}
