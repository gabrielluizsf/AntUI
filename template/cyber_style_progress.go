package template

import "github.com/gabrielluizsf/antui"

func (cyberStyle) Progress(win *antui.Window, x, y, w, h int, value float32) {
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)

	value = min(max(value, 0), 1)
	filled := int(float32(w)*value + 0.5)

	cv.FillRoundRect(x, y, w, h, h/2, cyberPalette.panelAct)
	cv.RoundRect(x, y, w, h, h/2, v.edge)
	if filled > 0 {
		glow(cv, x, y, filled, h, h/2, v.neon, 5, u)
		cv.FillRoundRect(x, y, filled, h, h/2, v.neon)
	}
}
