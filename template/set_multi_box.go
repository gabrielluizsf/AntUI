package template

import "github.com/gabrielluizsf/antui/template/css"

// setMultiBox gives a batch the geometry of the container's border box: where
// it sits, how wide it is, and the content box the columns are cut out of.
// The columns are solved here, before the children are measured, so a width in
// a percentage answers to the column it will be drawn in. pctBase is what a
// percentage of the container itself answers to, and a height of zero means the
// container is as tall as its columns turn out to be.
func (c *CSS) setMultiBox(b *multiBatch, st css.Style, x, y, w, h, pctBase int) {
	u := c.style.u()
	left := st.BorderWidth[3]*u + c.style.length(st, st.Padding[3], pctBase)
	right := st.BorderWidth[1]*u + c.style.length(st, st.Padding[1], pctBase)
	b.top = st.BorderWidth[0]*u + c.style.length(st, st.Padding[0], pctBase)
	b.bottom = st.BorderWidth[2]*u + c.style.length(st, st.Padding[2], pctBase)
	b.x, b.y, b.w, b.h = x, y, w, h
	b.originX, b.originY = x+left, y+b.top
	b.left, b.right = left, right
	b.contentW = max(w-left-right, 0)
	b.contentH = max(h-b.top-b.bottom, 0)
	b.gap = c.style.length(st, st.ColumnGap, b.contentW)
	b.solveColumns(st)
	// A height a stylesheet or a parent decided is one the flow can cut boxes
	// against. A height of zero is the container asking for whatever its columns
	// come to, which is the height nestedMultiCol's measure pass is about to
	// produce, and a parent that then hands that height back must not be read as
	// a decision — drawMultiContainer is the one that settles that.
	b.definite = h > 0
}
