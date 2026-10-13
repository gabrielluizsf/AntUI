package template

import "github.com/gabrielluizsf/antui/template/css"

// measureGridContainer lays a grid container out on its own, so its parent has
// a box to place: the width its own columns ask for, the height that width
// produces, and the narrowest of both. pctBase is what a percentage inside the
// container answers to, which is the width the parent would stretch it to.
func (c *CSS) measureGridContainer(st css.Style, draw func(*CSS), pctBase int) *gridNested {
	left, right, top, bottom := c.gridChrome(st, pctBase)
	height := c.gridDeclaredHeight(st)
	ml := c.style.length(st, st.Margin[css.SideLeft], pctBase)
	mr := c.style.length(st, st.Margin[css.SideRight], pctBase)
	width := max(pctBase-ml-mr, 0)
	if st.Has("width") && gridLengthDefinite(st.Width) {
		width = max(c.style.length(st, st.Width, pctBase)-left-right, 0)
	}

	b := &gridBatch{css: c, st: st}
	c.setGridBox(b, st, 0, 0, width+left+right, height, pctBase)
	// The silent pass of the whole page is already running, so the children are
	// collected here and the container is sized from them.
	c.gridWithin(b, draw)
	b.solve()

	// The narrowest the container can be, and the width it asks for: its own
	// columns say how wide, and the rows answer to whatever width that is.
	minW, widest := b.gridIntrinsicWidth()
	c.setGridBox(b, st, 0, 0, minW+left+right, height, pctBase)
	b.solve()
	minH := b.contentH + top + bottom
	if widest != b.contentW {
		c.setGridBox(b, st, 0, 0, widest+left+right, height, pctBase)
		b.solve()
	}
	return &gridNested{
		batch: b,
		nw:    b.contentW + left + right,
		nh:    b.contentH + top + bottom,
		minW:  minW + left + right,
		minH:  minH,
	}
}
