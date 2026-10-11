package template

import "github.com/gabrielluizsf/antui/template/css"

// baseStyle is the computed cascade style for a role and state, animations
// not yet folded in.
func (cs *cssStyle) baseStyle(role string, s State) css.Style {
	return cs.classes.GetStyleViewport(role, nil, css.StateWith(
		s.Hovered, s.Focused, s.Active || s.Pressed, s.On,
	), cs.viewport())
}
