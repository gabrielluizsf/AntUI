package template

import "github.com/gabrielluizsf/antui/template/css"

// closeMultiCol gives the batch its height — the tallest column, unless the
// container declared one and the flow cut against it — and its content box.
func (c *CSS) closeMultiCol(b *multiBatch, st css.Style) int {
	winH := c.win.Height()
	h := max(b.colH+b.top+b.bottom, c.clampDim(st, 0, "height", winH))
	if b.cuts {
		h = c.multiDeclaredHeight(st)
	}
	b.h = h
	b.contentH = max(h-b.top-b.bottom, 0)
	return h
}
