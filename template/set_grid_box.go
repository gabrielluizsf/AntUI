package template

import "github.com/gabrielluizsf/antui/template/css"

// setGridBox gives a batch the geometry of the container's border box on
// screen: where it sits, how big it is, and the content box its children are
// placed in. pctBase is what a percentage inside the box answers to, and a
// percentage row-gap is answered to what a declared height leaves, since a
// container sized by its content has no height to divide.
func (c *CSS) setGridBox(b *gridBatch, st css.Style, x, y, w, h, pctBase int) {
	left, right, top, bottom := c.gridChrome(st, pctBase)
	b.x, b.y, b.w, b.h = x, y, w, h
	b.originX, b.originY = x+left, y+top
	b.contentW = max(w-left-right, 0)
	b.contentH = max(h-top-bottom, 0)
	rowBase := b.contentH
	if st.Has("height") && gridLengthDefinite(st.Height) {
		rowBase = max(c.style.length(st, st.Height, c.win.Height())-top-bottom, 0)
	}
	b.rowGap = c.style.length(st, st.RowGap, rowBase)
	b.colGap = c.style.length(st, st.ColumnGap, b.contentW)
}
