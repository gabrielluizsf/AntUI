package template

import "github.com/gabrielluizsf/antui"

func (t *uiTemplate) Label(win *antui.Window, x, y int, text string) {
	t.style.Label(win, x, y, text)
}
