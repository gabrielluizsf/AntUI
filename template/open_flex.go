package template

import "github.com/gabrielluizsf/antui/template/css"

// openFlex closes the inline line for a flex container, works out its border
// box from the flow cursor and makes the batch the solver will fill.
func (c *CSS) openFlex(st css.Style) *flexBatch {
	u := c.style.u()
	winW := c.win.Width()
	ml := c.style.length(st, st.Margin[3], winW)
	mr := c.style.length(st, st.Margin[1], winW)
	mt := c.style.length(st, st.Margin[0], winW)

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

	bl := st.BorderWidth[3] * u
	br := st.BorderWidth[1] * u
	bt := st.BorderWidth[0] * u
	bb := st.BorderWidth[2] * u
	pl := c.style.length(st, st.Padding[3], winW)
	pr := c.style.length(st, st.Padding[1], winW)
	pt := c.style.length(st, st.Padding[0], winW)
	pb := c.style.length(st, st.Padding[2], winW)

	return &flexBatch{
		css:      c,
		st:       st,
		x:        x,
		y:        y,
		w:        w,
		originX:  x + bl + pl,
		contentW: max(w-br-bl-pr-pl, 0),
		frameW:   br + bl + pr + pl,
		frameH:   bb + bt + pb + pt,
		frameTop: bt + pt,
	}
}
