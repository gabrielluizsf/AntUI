package template

import "github.com/gabrielluizsf/antui/template/css"

// gridChrome is the room a container's own border and padding take out of its
// border box, with a percentage padding answered to pctBase.
func (c *CSS) gridChrome(st css.Style, pctBase int) (left, right, top, bottom int) {
	u := c.style.u()
	return st.BorderWidth[3]*u + c.style.length(st, st.Padding[3], pctBase),
		st.BorderWidth[1]*u + c.style.length(st, st.Padding[1], pctBase),
		st.BorderWidth[0]*u + c.style.length(st, st.Padding[0], pctBase),
		st.BorderWidth[2]*u + c.style.length(st, st.Padding[2], pctBase)
}
