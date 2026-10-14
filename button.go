package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// Button draws a button and reports whether it was clicked this frame.
func (win *Window) Button(area canvas.Area, label string) bool {
	x, y, w, h := area.X, area.Y, area.Width, area.Height
	id := widgetID("button", x, y, w, h, label)
	hovered := win.Hovered(area)
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
	win.cv.FillRoundRect(canvas.Round(area, win.theme.Radius), fill)

	if focused && !pressed {
		win.cv.RoundRect(canvas.Round(canvas.Area{X: x - 2, Y: y - 2, Width: w + 4, Height: h + 4}, win.theme.Radius+2), win.theme.Accent)
	}

	textY := y + (h-canvas.TextHeight())/2
	if pressed {
		textY++
	}
	win.cv.Text(x+(w-canvas.TextWidth(label))/2, textY, label, win.theme.TextOnAccent)
	return clicked
}
