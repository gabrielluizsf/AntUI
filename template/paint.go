package template

import "github.com/gabrielluizsf/antui/template/css"

// paint runs a widget's painter clipped to its box when overflow asks for it,
// then adds the scrollbar indicators for scroll/auto. A transform on the
// style records the whole pass on a layer and composites it back, so the box
// and its shadows rotate and scale together.
func (c *CSS) paint(st css.Style, role, label string, x, y, w, h int, draw func()) {
	if !visible(st) {
		return
	}
	c.style.transformed(st, x, y, w, h, func() {
		if interactive(st) && c.style.boxOn(st) {
			if ownSurface(role) {
				c.style.paintFrame(c.win, st, x, y, w, h)
			} else {
				c.style.paintBox(c.win, st, x, y, w, h)
			}
		}
		if st.Overflow == [2]uint8{css.OverflowVisible, css.OverflowVisible} {
			draw()
			return
		}
		u := c.style.u()
		nw, nh := c.style.natural(role, label, st, u)
		cv := c.win.Canvas()
		saved := cv.Clip
		cv.SetClip(x, y, w, h)
		draw()
		cv.Clip = saved
		c.style.paintScrollbars(c.win, st, x, y, w, h, nw, nh)
	})
}
