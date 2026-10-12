package template

import "github.com/gabrielluizsf/antui/template/css"

// multiColLayout is the container [CSS.MultiCol] starts: it measures the
// children, flows them into the columns and draws them at them.
func (c *CSS) multiColLayout(draw func(*CSS)) {
	st := c.style.baseStyle(css.RoleMultiCol, State{})
	if st.Display == css.DisplayNone {
		return
	}
	b := c.openMultiCol(st)
	c.multi, c.multicoling = b, true

	c.measure = true
	draw(c)
	c.measure = false
	b.flow()

	h := c.closeMultiCol(b, st)

	if c.style.boxOn(st) && visible(st) {
		c.style.paintBox(c.win, st, b.x, b.y, b.w, h)
	}
	c.style.paintColumnRules(c.win, b)

	for pass := range b.passes {
		b.pass = pass
		c.multiWithin(b, draw)
	}

	c.multicoling = false
	c.multi = nil
	c.cursorY = b.y + h
	c.lastMB = c.style.length(st, st.Margin[2], c.win.Width())
}
