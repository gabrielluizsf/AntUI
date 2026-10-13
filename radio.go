package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// Radio draws one option of a group. Clicking it sets value to option and
// reports true; clicking the already-selected option changes nothing.
func (win *Window) Radio(x, y int, label string, value *int, option int) bool {
	if value == nil {
		return false
	}
	const size = 18
	w := size + 8 + canvas.TextWidth(label)
	id := widgetID("radio", x, y, size, option, label)
	hovered := win.Hovered(x, y, w, size)
	clicked := win.widgetClick(id, hovered)
	focused := win.uiFocus == id
	win.uiTab = append(win.uiTab, id)
	if win.WidgetKeyActivate(id) {
		clicked = true
	}
	selected := *value == option

	if clicked && !selected {
		*value = option
		selected = true
	} else {
		clicked = false
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

	win.cv.Text(x+size+8, y+(size-canvas.TextHeight())/2, label, win.theme.Text)
	return clicked
}
