package template

import "github.com/gabrielluizsf/antui/template/css"

// openFlex closes the inline line for a flex container, works out its border
// box from the flow cursor and makes the batch the solver will fill.
func (c *CSS) openFlex(st css.Style) *flexBatch {
	u := c.style.u()
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

	bl := st.BorderWidth[css.SideLeft] * u
	br := st.BorderWidth[css.SideRight] * u
	bt := st.BorderWidth[css.SideTop] * u
	bb := st.BorderWidth[css.SideBottom] * u
	pl := c.style.length(st, st.Padding[css.SideLeft], winW)
	pr := c.style.length(st, st.Padding[css.SideRight], winW)
	pt := c.style.length(st, st.Padding[css.SideTop], winW)
	pb := c.style.length(st, st.Padding[css.SideBottom], winW)

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
