package template

import "github.com/gabrielluizsf/antui/template/css"

// multiItemSized collects one child the way the stylesheet asks for it, from a
// size the child wants: a widget's own natural box, or the box a nested
// container asked for. The properties a multi-column flow reads are the same
// either way, so a container drawn inside this one is placed as a real child.
func (c *CSS) multiItemSized(st css.Style, nw, nh int) multiItem {
	b := c.multi
	return multiItem{
		st: st,
		nw: nw, nh: nh, w: nw, h: nh,
		ml:   c.style.length(st, st.Margin[css.SideLeft], b.colW),
		mr:   c.style.length(st, st.Margin[css.SideRight], b.colW),
		mt:   c.style.length(st, st.Margin[css.SideTop], b.colW),
		mb:   c.style.length(st, st.Margin[css.SideBottom], b.colW),
		keep: st.BreakInside == css.BreakAvoid,
	}
}
