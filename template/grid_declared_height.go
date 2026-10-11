package template

import "github.com/gabrielluizsf/antui/template/css"

// gridDeclaredHeight is the border-box height the stylesheet gives a container,
// clamped by its own min and max, and 0 when it declares none. It bounds the
// content box before the rows are sized, which is what gives them a space to
// fill and a definite axis to divide.
func (c *CSS) gridDeclaredHeight(st css.Style) int {
	if !st.Has("height") || !gridLengthDefinite(st.Height) {
		return 0
	}
	winH := c.win.Height()
	return max(c.clampDim(st, c.style.length(st, st.Height, winH), "height", winH), 0)
}
