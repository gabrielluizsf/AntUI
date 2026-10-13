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

// styleForClass is the style of a widget drawn from another widget's role
// that carries a class of the template's own — the calendar the date box
// drops, which has to stay styleable apart from the box while the element's
// own rules keep reaching it. It paints as the cascade answers it, with no
// animation of the box's to inherit.
func (cs *cssStyle) styleForClass(role, class string, s State) css.Style {
	return cs.baseStyleClass(role, []string{class}, s)
}
