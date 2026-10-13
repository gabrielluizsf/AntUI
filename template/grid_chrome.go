package template

import "github.com/gabrielluizsf/antui/template/css"

// gridChrome is the room a container's own border and padding take out of its
// border box, with a percentage padding answered to pctBase.
func (c *CSS) gridChrome(st css.Style, pctBase int) (left, right, top, bottom int) {
	u := c.style.u()
	return st.BorderWidth[css.SideLeft]*u + c.style.length(st, st.Padding[css.SideLeft], pctBase),
		st.BorderWidth[css.SideRight]*u + c.style.length(st, st.Padding[css.SideRight], pctBase),
		st.BorderWidth[css.SideTop]*u + c.style.length(st, st.Padding[css.SideTop], pctBase),
		st.BorderWidth[css.SideBottom]*u + c.style.length(st, st.Padding[css.SideBottom], pctBase)
}
