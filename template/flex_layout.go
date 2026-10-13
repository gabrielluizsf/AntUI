package template

import "github.com/gabrielluizsf/antui/template/css"

// flexLayout is the container [CSS.Flex] starts: it measures the children,
// solves their boxes and draws them at them.
func (c *CSS) flexLayout(draw func(*CSS)) {
	st := c.style.baseStyle(css.RoleFlex, State{})
	if st.Display == css.DisplayNone {
		return
	}
	b := c.openFlex(st)
	c.flex = b
	c.flexing = true

	c.measure = true
	draw(c)
	c.measure = false

	// A container that declares its height bounds the content box up front,
	// so the solver has a definite main or cross size to distribute against.
	if st.Has("height") {
		winH := c.win.Height()
		hh := c.clampDim(st, c.style.length(st, st.Height, winH), "height", winH)
		b.contentH = max(hh-b.frameH, 0)
	}
	b.solve()

	h := c.closeFlex(b, st)

	if c.style.boxOn(st) && visible(st) {
		c.style.paintBox(c.win, st, b.x, b.y, b.w, h)
	}

	b.next = 0
	draw(c)

	c.flexing = false
	c.flex = nil

	c.cursorY = b.y + h
	c.lastMB = c.style.length(st, st.Margin[css.SideBottom], c.win.Width())
}
