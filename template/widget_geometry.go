package template

import "github.com/gabrielluizsf/antui"

// WidgetGeometry is the common geometry for all widget painting methods.
type WidgetGeometry struct {
	Win   *antui.Window
	State State
	X, Y  int
	W, H  int
}