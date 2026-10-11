package template

import "github.com/gabrielluizsf/antui/template/css"

// multiDeclaredHeight is the border-box height the stylesheet gives a
// container, clamped by its own min and max, and 0 when it declares none.
func (c *CSS) multiDeclaredHeight(st css.Style) int {
	if !st.Has("height") || !gridLengthDefinite(st.Height) {
		return 0
	}
	winH := c.win.Height()
	return max(c.clampDim(st, c.style.length(st, st.Height, winH), "height", winH), 0)
}
