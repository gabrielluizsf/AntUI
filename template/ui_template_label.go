package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

func (t *uiTemplate) Label(win *antui.Window, area canvas.Area, text string) {
	x, y := area.X, area.Y
	t.style.Label(LabelParams{
		Win: win,
		X:   x, Y: y,
		Text: text,
	})
}
