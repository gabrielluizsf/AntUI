package template

import "github.com/gabrielluizsf/antui/template/css"

// openMultiCol closes the inline line for a multi-column container, works out
// its border box from the flow cursor and makes the batch the flow will fill.
func (c *CSS) openMultiCol(st css.Style) *multiBatch {
	winW := c.win.Width()
	ml := c.style.length(st, st.Margin[css.SideLeft], winW)
	mr := c.style.length(st, st.Margin[css.SideRight], winW)
	mt := c.style.length(st, st.Margin[css.SideTop], winW)

	// The container is a block: it closes the inline line, collapses the
	// sibling margins and rests on the flow cursor.
	c.cursorY += c.lineH
	c.lineH = 0
	c.lineX = 0
	c.cursorY += max(c.lastMB, mt)
	y := c.cursorY

	w := max(winW-ml-mr, 0)
	if st.Has("width") {
		w = c.style.length(st, st.Width, winW)
	}
	w = c.clampDim(st, w, "width", winW)
	x := ml + (winW-ml-mr-w)/2

	b := &multiBatch{css: c, st: st}
	c.setMultiBox(b, st, x, y, w, c.multiDeclaredHeight(st), winW)
	return b
}
