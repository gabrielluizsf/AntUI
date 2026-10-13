package template

import "github.com/gabrielluizsf/antui/template/css"

// tableChrome is the border and padding the table's own style adds around the
// cells, which a percentage inside it answers to.
func (c *CSS) tableChrome(st css.Style, pctBase int) (left, right, top, bottom int) {
	u := c.style.u()
	return st.BorderWidth[css.SideLeft]*u + c.style.length(st, st.Padding[css.SideLeft], pctBase),
		st.BorderWidth[css.SideRight]*u + c.style.length(st, st.Padding[css.SideRight], pctBase),
		st.BorderWidth[css.SideTop]*u + c.style.length(st, st.Padding[css.SideTop], pctBase),
		st.BorderWidth[css.SideBottom]*u + c.style.length(st, st.Padding[css.SideBottom], pctBase)
}
