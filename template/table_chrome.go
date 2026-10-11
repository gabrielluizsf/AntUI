package template

import "github.com/gabrielluizsf/antui/template/css"

// tableChrome is the border and padding the table's own style adds around the
// cells, which a percentage inside it answers to.
func (c *CSS) tableChrome(st css.Style, pctBase int) (left, right, top, bottom int) {
	u := c.style.u()
	return st.BorderWidth[3]*u + c.style.length(st, st.Padding[3], pctBase),
		st.BorderWidth[1]*u + c.style.length(st, st.Padding[1], pctBase),
		st.BorderWidth[0]*u + c.style.length(st, st.Padding[0], pctBase),
		st.BorderWidth[2]*u + c.style.length(st, st.Padding[2], pctBase)
}
