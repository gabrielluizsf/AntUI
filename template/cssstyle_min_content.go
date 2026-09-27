package template

// The narrowest a widget can be: what a grid track asking for min-content gets,
// measured the way the text model wraps — down to the widest word that cannot
// be broken.

import (
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// minContent is the size a widget takes when nothing else is asked of it: the
// widest word it cannot break, plus whatever it draws either way.
func (cs *cssStyle) minContent(role, label string, st css.Style, u int) (w, h int) {
	return cs.sized(role, st, u, cs.widestWord(st, label))
}

// widestWord is the width of the longest run of text the widget cannot break
// across two lines, measured with the style's own typography. Whitespace alone
// is no width at all.
func (cs *cssStyle) widestWord(st css.Style, label string) int {
	if label == "" {
		return 0
	}
	text := transformText(label, st.TextTransform)
	if !strings.ContainsAny(text, " \t\n\r") {
		return cs.measure(st, text)
	}
	widest := 0
	for _, word := range strings.Fields(text) {
		widest = max(widest, cs.measure(st, word))
	}
	return widest
}

// sized is the shape of a widget built from the width of its text, so the
// natural reading and the min-content reading of one role cannot drift apart.
func (cs *cssStyle) sized(role string, st css.Style, u, labelW int) (w, h int) {
	var bx, by int
	if st.BorderOn() {
		bx = (st.BorderWidth[1] + st.BorderWidth[3]) * u
		by = (st.BorderWidth[0] + st.BorderWidth[2]) * u
	}
	ph := textHeight(u)
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
