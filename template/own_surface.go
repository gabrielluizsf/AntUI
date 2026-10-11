package template

import "github.com/gabrielluizsf/antui/template/css"

// ownSurface reports whether the widget's own painter draws the background its
// style declares. The three toggles are controls of a size of their own, drawn
// inside a cell the solver gave them: a checkbox is a box of eighteen pixels
// with its label beside it, not a surface running under the label, and a switch
// is a track of thirty-six. Their border and shadow stay the cell's, which is
// what a rule drawn under a row of a table is.
func ownSurface(role string) bool {
	switch role {
	case css.RoleCheckbox, css.RoleRadio, css.RoleSwitch:
		return true
	}
	return false
}
