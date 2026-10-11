package template

import "github.com/gabrielluizsf/antui/template/css"

// placeOutOfFlow positions an absolute or fixed box from its inset edges:
// left/top when given, otherwise right/bottom measured from the far edge of
// the window (the viewport). With neither, the box rests at the window origin.
func (c *CSS) placeOutOfFlow(st css.Style, w, h, winW, winH int) (int, int) {
	x, y := 0, 0
	switch {
	case st.Has("left"):
		x = c.style.length(st, st.Left, winW)
	case st.Has("right"):
		x = winW - w - c.style.length(st, st.Right, winW)
	}
	switch {
	case st.Has("top"):
		y = c.style.length(st, st.Top, winH)
	case st.Has("bottom"):
		y = winH - h - c.style.length(st, st.Bottom, winH)
	}
	return x, y
}
