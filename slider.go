package antui

// Slider draws a horizontal slider and reports whether the value changed
// this frame.
func (win *Window) Slider(x, y, w, h int, value *float32, minValue, maxValue float32) bool {
	if value == nil || w <= 0 || h <= 0 || maxValue <= minValue {
		return false
	}
	id := widgetID("slider", x, y, w, h, "")
	hovered := win.Hovered(x, y, w, h)
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
