package antui

import "github.com/gabrielluizsf/antui/canvas"

// Slider draws a horizontal slider and reports whether the value changed
// this frame.
func (win *Window) Slider(area canvas.Area, params SliderParams) bool {
	if params.Value == nil || area.Width <= 0 || area.Height <= 0 || params.Max <= params.Min {
		return false
	}
	value, minValue, maxValue := params.Value, params.Min, params.Max
	x, y, w, h := area.X, area.Y, area.Width, area.Height
	id := widgetID("slider", x, y, w, h, "")
	hovered := win.Hovered(area)
	focused := win.uiFocus == id
	win.uiTab = append(win.uiTab, id)

	if hovered {
		win.uiHot = id
	}
	if hovered && win.MousePressed(MouseLeft) {
		win.uiActive = id
		win.uiFocus = 0
	}
	if win.uiActive == id && !win.MouseDown(MouseLeft) {
		win.uiActive = 0
	}

	s := sliderState{
		x: x, y: y, w: w, h: h,
		id: id, hovered: hovered, focused: focused,
		knobRadius: max(h/2, 6),
		minValue:   minValue,
		maxValue:   maxValue,
	}
	s.usable = max(w-2*s.knobRadius, 1)
	s.span = maxValue - minValue

	changed := win.sliderInput(s, value)

	*value = min(max(*value, minValue), maxValue)
	win.sliderDraw(s, *value)
	return changed
}
