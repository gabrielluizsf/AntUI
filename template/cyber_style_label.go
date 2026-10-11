package template

import "github.com/gabrielluizsf/antui"

func (cyberStyle) Label(win *antui.Window, x, y int, text string) {
	win.Canvas().TextScaled(x, y, text, vividFor(win).text, Scale(win))
}
