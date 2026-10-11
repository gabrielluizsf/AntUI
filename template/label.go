package template

import "github.com/gabrielluizsf/antui"

func (builtinStyle) Label(win *antui.Window, x, y int, text string) {
	win.Canvas().TextScaled(x, y, text, win.Theme().Text, Scale(win))
}
