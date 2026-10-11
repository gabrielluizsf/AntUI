package template

import "github.com/gabrielluizsf/antui/template/css"

// layout measures and places one widget in the flow, opening the widget's
// timeline entry and hosting its style on it. ok is false for a widget the
// stylesheet hides (display:none), which draws nothing and hears nothing.
// An absolutely- or fixed-positioned widget is placed by its inset edges and
// does not advance the flow; a relative or sticky one keeps its slot and only
// shifts; every other widget flows down, or along a line when inline. Inside
// a [CSS.Flex] or [CSS.Grid] block a widget is an item of that container
// instead: the silent pass records it and reports ok=false, and the draw pass
// hands back the box the solver placed it in.
func (c *CSS) layout(role, label string) (x, y, w, h int, st css.Style, ok bool) {
	if c.flexing {
		return c.layoutFlex(role, label)
	}
	if c.griding {
		return c.layoutGrid(role, label)
	}
	if c.multicoling {
		return c.layoutMulti(role, label)
	}
	if c.tabling {
		return c.layoutTable(role, label)
	}
	e := c.style.beginWidget(role, label)
	return c.layoutBox(e, role, label, e.state)
}
