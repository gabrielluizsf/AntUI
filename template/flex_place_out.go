package template

import "github.com/gabrielluizsf/antui/template/css"

// flexPlaceOut positions an out-of-flow child by its inset edges against the
// container's content box, left/top when given and right/bottom from the far
// edge otherwise.
func (c *CSS) flexPlaceOut(st css.Style, it flexItem) (x, y, w, h int, os css.Style, ok bool) {
	b := c.flex
	x, y = b.originX, b.originY
	if it.inLeft >= 0 {
		x += it.inLeft
	} else if it.inRight >= 0 {
		x += b.contentW - it.w - it.inRight
	}
	if it.inTop >= 0 {
		y += it.inTop
	} else if it.inBottom >= 0 {
		y += b.contentH - it.h - it.inBottom
	}
	return x, y, it.w, it.h, st, true
}
