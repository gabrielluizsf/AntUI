package template

import "github.com/gabrielluizsf/antui/template/css"

// styleFor resolves the animated style for a widget role, folding the
// interaction state into the pseudo-classes. While a widget is being laid
// out, styleFor drives its timeline entry; otherwise the plain cascade
// answer is returned, which keeps styles the template reads for itself free
// of any frame's animation.
func (cs *cssStyle) styleFor(role string, s State) css.Style {
	if cs.cur != nil && cs.cur.role == role {
		return cs.entryStyle(cs.cur, role, s)
	}
	return cs.baseStyle(role, s)
}
