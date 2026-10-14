package template

import "github.com/gabrielluizsf/antui/template/css"

// baseStyle is the computed cascade style for a role and state, animations
// not yet folded in.
func (cs *cssStyle) baseStyle(role string, s State) css.Style {
	return cs.baseStyleClass(role, nil, s)
}

// baseStyleClass is baseStyle with classes of the template's own on the
// widget, which the class table does not hold for its kind.
func (cs *cssStyle) baseStyleClass(role string, extra []string, s State) css.Style {
	return cs.classes.GetStyleViewport(css.StyleQuery{
		Tag:     role,
		Classes: extra,
		State: css.StateWith(css.StateFlags{
			Hovered: s.Hovered,
			Focused: s.Focused,
			Active:  s.Active || s.Pressed,
			Checked: s.On,
		}),
		Viewport: cs.viewport(),
	})
}
