package template

import "github.com/gabrielluizsf/antui/template/css"

// sized is the shape of a widget built from the width of its text, so the
// natural reading and the min-content reading of one role cannot drift apart.
func (cs *cssStyle) sized(role string, st css.Style, u, labelW int) (w, h int) {
	var bx, by int
	if st.BorderOn() {
		bx = (st.BorderWidth[css.SideRight] + st.BorderWidth[css.SideLeft]) * u
		by = (st.BorderWidth[css.SideTop] + st.BorderWidth[css.SideBottom]) * u
	}
	ph := cs.textHeight(st)
	switch role {
	case css.RoleLabel:
		return labelW, ph
	case css.RoleButton:
		return labelW + 2*cs.padding(st, 3) + bx,
			ph + cs.padding(st, 2) + cs.padding(st, 0) + by
	case css.RoleCheckbox, css.RoleRadio:
		return boxSize*u + labelGap*u + labelW, boxSize * u
	case css.RoleSlider:
		return controlW * u, sliderHeight * u
	case css.RoleInput:
		return max(labelW+2*cs.padding(st, 3)+bx, controlW*u),
			ph + cs.padding(st, 2) + cs.padding(st, 0) + by
	case css.RoleSelect:
		return max(labelW+2*cs.padding(st, 3)+caretW*u+bx, controlW*u),
			ph + cs.padding(st, 2) + cs.padding(st, 0) + by
	case css.RoleTextArea:
		return max(controlW*u, textAreaMinCols*u), sliderHeight * u
	case css.RoleSwitch:
		return switchTrackW*u + labelGap*u + labelW, switchTrackH * u
	case css.RoleProgress:
		return controlW * u, progressHeight * u
	case css.RoleDatePicker:
		return max(labelW+2*cs.padding(st, 3)+caretW*u+bx, controlW*u),
			ph + cs.padding(st, 2) + cs.padding(st, 0) + by
	}
	return labelW, ph
}
