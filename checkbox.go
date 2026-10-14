package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// Checkbox draws a labelled checkbox, flips value when clicked, and reports
// whether it changed this frame.
func (win *Window) Checkbox(x, y int, label string, value *bool) bool {
	if value == nil {
		return false
	}
	const box = 18
	w := box + 8 + canvas.TextWidth(label)
	id := widgetID("checkbox", x, y, box, box, label)
	hovered := win.Hovered(canvas.Area{X: x, Y: y, Width: w, Height: box})
	clicked := win.widgetClick(id, hovered)
	focused := win.uiFocus == id
	win.uiTab = append(win.uiTab, id)
	if win.WidgetKeyActivate(id) {
		clicked = true
	}
	if clicked {
		*value = !*value
	}

	fill := win.theme.Surface
	switch {
	case *value:
		fill = win.theme.Accent
	case hovered:
		fill = win.theme.SurfaceHover
	}
	border := win.theme.Border
	if *value {
		border = win.theme.Accent
	}
	win.cv.FillRoundRect(canvas.Round(canvas.Area{X: x, Y: y, Width: box, Height: box}, 4), fill)
	win.cv.RoundRect(canvas.Round(canvas.Area{X: x, Y: y, Width: box, Height: box}, 4), border)

	if focused {
		win.cv.RoundRect(canvas.Round(canvas.Area{X: x - 2, Y: y - 2, Width: box + 4, Height: box + 4}, 6), win.theme.Accent)
	}

	if *value {
		mark := win.theme.TextOnAccent
		win.cv.Line(x+4, y+9, x+7, y+13, mark)
		win.cv.Line(x+5, y+9, x+8, y+13, mark)
		win.cv.Line(x+7, y+13, x+14, y+5, mark)
		win.cv.Line(x+8, y+13, x+15, y+5, mark)
	}

	win.cv.Text(x+box+8, y+(box-canvas.TextHeight())/2, label, win.theme.Text)
	return clicked
}
