package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// Button draws a button and reports whether it was clicked this frame.
func (win *Window) Button(x, y, w, h int, label string) bool {
	id := widgetID("button", x, y, w, h, label)
	hovered := win.Hovered(x, y, w, h)
	clicked := win.widgetClick(id, hovered)
	focused := win.uiFocus == id
	win.uiTab = append(win.uiTab, id)
	keyed := win.WidgetKeyActivate(id)
	if keyed {
		clicked = true
	}
	pressed := (win.uiActive == id && hovered) || keyed

	fill := win.theme.Accent
	switch {
	case pressed:
		fill = canvas.Shade(win.theme.Accent, -0.18)
	case hovered:
		fill = win.theme.AccentHover
	}
	win.cv.FillRoundRect(x, y, w, h, win.theme.Radius, fill)

	if focused && !pressed {
		win.cv.RoundRect(x-2, y-2, w+4, h+4, win.theme.Radius+2, win.theme.Accent)
	}

	textY := y + (h-canvas.TextHeight())/2
	if pressed {
		textY++
	}
	win.cv.Text(x+(w-canvas.TextWidth(label))/2, textY, label, win.theme.TextOnAccent)
	return clicked
}
