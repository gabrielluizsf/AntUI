package template

import "github.com/gabrielluizsf/antui/template/css"

// shiftInset moves a relative or sticky box away from its flow position by its
// inset edges, leaving the slot it occupies untouched.
func (c *CSS) shiftInset(st css.Style, x, y, winW, winH int) (int, int) {
	if st.Position != css.PositionRelative && st.Position != css.PositionSticky {
		return x, y
	}
	switch {
	case st.Has("left"):
		x += c.style.length(st, st.Left, winW)
	case st.Has("right"):
		x -= c.style.length(st, st.Right, winW)
	}
	switch {
	case st.Has("top"):
		y += c.style.length(st, st.Top, winH)
	case st.Has("bottom"):
		y -= c.style.length(st, st.Bottom, winH)
	}
	return x, y
}
