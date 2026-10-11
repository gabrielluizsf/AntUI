package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

func (t *uiTemplate) Background(win *antui.Window) canvas.Color {
	return t.style.Background(win)
}
