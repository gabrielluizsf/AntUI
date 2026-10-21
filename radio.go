package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// Radio draws one option of a group. Clicking an active option clears the
// group by storing NoOption, so the group can have no selection.
func (win *Window) Radio(x, y int, params RadioParams) bool {
	if params.Value == nil {
		return false
	}
	label, value, option := params.Label, params.Value, params.Option
	const size = 18
	w := size + 8 + canvas.TextWidth(label)
	id := widgetID("radio", x, y, size, option, label)
	hovered := win.Hovered(canvas.Area{X: x, Y: y, Width: w, Height: size})
	clicked := win.widgetClick(id, hovered)
	focused := win.uiFocus == id
	win.uiTab = append(win.uiTab, id)
	if win.WidgetKeyActivate(id) {
		clicked = true
	}
	selected := *value == option

	if clicked {
		if selected {
			*value = NoOption
			selected = false
		} else {
			*value = option
			selected = true
		}
	}

	surface := win.theme.Surface
	if hovered {
		surface = win.theme.SurfaceHover
	}
	border := win.theme.Border
	if selected {
		border = win.theme.Accent
	}
	win.cv.FillCircle(canvas.Circle{X: x + size/2, Y: y + size/2, R: size / 2}, surface)
	win.cv.Circle(canvas.Circle{X: x + size/2, Y: y + size/2, R: size / 2}, border)
	if focused {
		win.cv.Circle(canvas.Circle{X: x + size/2, Y: y + size/2, R: size/2 + 3}, win.theme.Accent)
	}
	if selected {
		win.cv.FillCircle(canvas.Circle{X: x + size/2, Y: y + size/2, R: size/2 - 4}, win.theme.Accent)
	}

	win.cv.Text(canvas.Point{X: float64(x + size + 8), Y: float64(y + (size-canvas.TextHeight())/2)}, label, win.theme.Text)
	return clicked
}
