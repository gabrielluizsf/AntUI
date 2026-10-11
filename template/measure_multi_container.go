package template

import "github.com/gabrielluizsf/antui/template/css"

// measureMultiContainer lays a multi-column container out on its own, so its
// parent has a box to place: the width the parent would stretch it to, and the
// height that width produces. pctBase is what a percentage of the container
// answers to. The children are collected once and re-flowed at the narrower
// width a parent that shrinks its items may end up giving.
func (c *CSS) measureMultiContainer(st css.Style, draw func(*CSS), pctBase int) *multiNested {
	ml := c.style.length(st, st.Margin[3], pctBase)
	mr := c.style.length(st, st.Margin[1], pctBase)
	width := max(pctBase-ml-mr, 0)
	if st.Has("width") {
		width = c.style.length(st, st.Width, pctBase)
	}
	width = c.clampDim(st, width, "width", pctBase)

	b := &multiBatch{css: c, st: st}
	c.setMultiBox(b, st, 0, 0, width, c.multiDeclaredHeight(st), pctBase)
	// The silent pass of the whole page is already running, so the children are
	// collected here and the container is sized from them.
	c.multiWithin(b, draw)
	b.flow()
	return &multiNested{batch: b, nw: b.ask(), nh: b.h}
}
