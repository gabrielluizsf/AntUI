package template

import "github.com/gabrielluizsf/antui/template/css"

// gridLayout is the container [CSS.Grid] starts: it measures the children,
// solves their cells and draws them at them.
func (c *CSS) gridLayout(draw func(*CSS)) {
	st := c.style.baseStyle(css.RoleGrid, State{})
	if st.Display == css.DisplayNone {
		return
	}
	winW := c.win.Width()
	ml := c.style.length(st, st.Margin[3], winW)
	mr := c.style.length(st, st.Margin[1], winW)
	mt := c.style.length(st, st.Margin[0], winW)
	mb := c.style.length(st, st.Margin[2], winW)

	c.cursorY += c.lineH
	c.lineH = 0
	c.lineX = 0
	c.cursorY += max(c.lastMB, mt)
	y := c.cursorY

	w := max(winW-ml-mr, 0)
	if st.Has("width") && gridLengthDefinite(st.Width) {
		w = c.style.length(st, st.Width, winW)
	}
	w = c.clampDim(st, w, "width", winW)
	x := ml + (winW-ml-mr-w)/2

	b := &gridBatch{css: c, st: st}
	c.setGridBox(b, st, x, y, w, c.gridDeclaredHeight(st), winW)
	c.grid, c.griding = b, true
	c.measure = true
	draw(c)
	c.measure = false
	b.solve()

	_, _, top, bottom := c.gridChrome(st, winW)
	h := c.gridHeight(st, b, top, bottom)
	c.setGridBox(b, st, x, y, w, h, winW)

	if c.style.boxOn(st) && visible(st) {
		c.style.paintBox(c.win, st, x, y, w, h)
	}

	c.gridWithin(b, draw)

	c.grid = nil
	c.griding = false
	c.cursorY = y + h
	c.lastMB = mb
}
