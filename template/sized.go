package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// sized is the shape of a widget built from the width of its text, so the
// natural reading and the min-content reading of one role cannot drift apart.
func (cs *cssStyle) sized(role string, st css.Style, u, labelW int) (w, h int) {
	var bx, by int
	if st.BorderOn() {
		bx = (st.BorderWidth[1] + st.BorderWidth[3]) * u
		by = (st.BorderWidth[0] + st.BorderWidth[2]) * u
	}
	ph := cs.textHeight(st)
	switch role {
	case css.RoleLabel:
		return labelW, ph
	case css.RoleButton:
		return labelW + 2*cs.padding(st, 3) + bx,
			ph + cs.padding(st, 2) + cs.padding(st, 0) + by
	case css.RoleCheckbox, css.RoleRadio:
		return 18*u + 8*u + labelW, 18 * u
	case css.RoleSlider:
		return 24 * canvas.FontWidth * u, 3 * canvas.FontHeight * u
	case css.RoleInput:
		return max(labelW+2*cs.padding(st, 3)+bx, 24*canvas.FontWidth*u),
			ph + cs.padding(st, 2) + cs.padding(st, 0) + by
	case css.RoleSelect:
		return max(labelW+2*cs.padding(st, 3)+16*u+bx, 24*canvas.FontWidth*u),
			ph + cs.padding(st, 2) + cs.padding(st, 0) + by
	case css.RoleTextArea:
		return max(24*canvas.FontWidth*u, 32*canvas.FontWidth*u), 3 * canvas.FontHeight * u
	case css.RoleSwitch:
		return 36*u + 8*u + labelW, 20 * u
	case css.RoleProgress:
		return 24 * canvas.FontWidth * u, 2 * canvas.FontHeight * u
	case css.RoleDatePicker:
		return max(labelW+2*cs.padding(st, 3)+16*u+bx, 24*canvas.FontWidth*u),
			ph + cs.padding(st, 2) + cs.padding(st, 0) + by
	}
	return labelW, ph
}
