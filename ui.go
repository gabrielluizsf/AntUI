package antui

import (
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// The ready-made widgets are immediate mode: there is no widget object and
// nothing to keep in sync. Each call draws the control and reports what the
// user did to it this frame.

// widgetID hashes a widget's identity with FNV-1a.
func widgetID(kind string, x, y, w, h int, label string) uint32 {
	const (
		offset = 2166136261
		prime  = 16777619
	)
	hash := uint32(offset)
	mix := func(b byte) {
		hash ^= uint32(b)
		hash *= prime
	}
	for _, v := range [4]int{x, y, w, h} {
		for shift := range 4 {
			mix(byte(uint32(v) >> (shift * 8)))
		}
	}
	for i := range len(kind) {
		mix(kind[i])
	}
	for i := range len(label) {
		mix(label[i])
	}
	if hash == 0 {
		return 1
	}
	return hash
}

// Hovered reports whether the pointer is inside a rectangle.
func (win *Window) Hovered(x, y, w, h int) bool {
	return win.mouseX >= x && win.mouseX < x+w &&
		win.mouseY >= y && win.mouseY < y+h
}

// Background paints the window in the theme's background colour.
func (win *Window) Background() { win.cv.Clear(win.theme.Background) }

// Label draws plain text in the theme's text colour.
func (win *Window) Label(x, y int, text string) {
	win.cv.Text(x, y, text, win.theme.Text)
}

// widgetClick is the press-and-release logic every clickable widget shares.
func (win *Window) widgetClick(id uint32, hovered bool) bool {
	if hovered {
		win.uiHot = id
	}
	switch {
	case win.uiActive == id:
		if win.MouseReleased(MouseLeft) {
			win.uiActive = 0
			return hovered
		}
	case hovered && win.MousePressed(MouseLeft):
		win.uiActive = id
		win.uiFocus = 0
		if win.MouseReleased(MouseLeft) {
			win.uiActive = 0
			return true
		}
	}
	return false
}

// Clicked reports whether a rectangle was clicked or tapped this frame,
// drawing nothing at all.
func (win *Window) Clicked(x, y, w, h int) bool {
	return win.widgetClick(widgetID("area", x, y, w, h, ""),
		win.Hovered(x, y, w, h))
}

// ClickedID is [Window.Clicked] for a control that moves: the name decides
// what "the same control" means, rather than where it happens to be.
func (win *Window) ClickedID(name string, x, y, w, h int) bool {
	return win.widgetClick(widgetID("area", 0, 0, 0, 0, name),
		win.Hovered(x, y, w, h))
}

// Button draws a button and reports whether it was clicked this frame.
func (win *Window) Button(x, y, w, h int, label string) bool {
	id := widgetID("button", x, y, w, h, label)
	hovered := win.Hovered(x, y, w, h)
	clicked := win.widgetClick(id, hovered)
	pressed := win.uiActive == id && hovered

	fill := win.theme.Accent
	switch {
	case pressed:
		fill = canvas.Shade(win.theme.Accent, -0.18)
	case hovered:
		fill = win.theme.AccentHover
	}
	win.cv.FillRoundRect(x, y, w, h, win.theme.Radius, fill)

	textY := y + (h-canvas.FontHeight)/2
	if pressed {
		textY++
	}
	win.cv.Text(x+(w-canvas.TextWidth(label))/2, textY, label, win.theme.TextOnAccent)
	return clicked
}

// Checkbox draws a labelled checkbox, flips value when clicked, and reports
// whether it changed this frame.
func (win *Window) Checkbox(x, y int, label string, value *bool) bool {
	if value == nil {
		return false
	}
	const box = 18
	w := box + 8 + canvas.TextWidth(label)
	id := widgetID("checkbox", x, y, box, box, label)
	hovered := win.Hovered(x, y, w, box)
	clicked := win.widgetClick(id, hovered)
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
	win.cv.FillRoundRect(x, y, box, box, 4, fill)
	win.cv.RoundRect(x, y, box, box, 4, border)

	if *value {
		mark := win.theme.TextOnAccent
		win.cv.Line(x+4, y+9, x+7, y+13, mark)
		win.cv.Line(x+5, y+9, x+8, y+13, mark)
		win.cv.Line(x+7, y+13, x+14, y+5, mark)
		win.cv.Line(x+8, y+13, x+15, y+5, mark)
	}

	win.cv.Text(x+box+8, y+(box-canvas.FontHeight)/2, label, win.theme.Text)
	return clicked
}

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
	win.cv.FillCircle(x+size/2, y+size/2, size/2, surface)
	win.cv.Circle(x+size/2, y+size/2, size/2, border)
	if selected {
		win.cv.FillCircle(x+size/2, y+size/2, size/2-4, win.theme.Accent)
	}

	win.cv.Text(x+size+8, y+(size-canvas.FontHeight)/2, label, win.theme.Text)
	return clicked
}

// Slider draws a horizontal slider and reports whether the value changed
// this frame.
func (win *Window) Slider(x, y, w, h int, value *float32, minValue, maxValue float32) bool {
	if value == nil || w <= 0 || h <= 0 || maxValue <= minValue {
		return false
	}
	id := widgetID("slider", x, y, w, h, "")
	hovered := win.Hovered(x, y, w, h)

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

	knobRadius := max(h/2, 6)
	usable := max(w-2*knobRadius, 1)
	span := maxValue - minValue

	changed := false
	if win.uiActive == id {
		position := float32(win.mouseX-(x+knobRadius)) / float32(usable)
		position = min(max(position, 0), 1)
		if wanted := minValue + position*span; wanted != *value {
			*value = wanted
			changed = true
		}
	}

	*value = min(max(*value, minValue), maxValue)
	t := (*value - minValue) / span

	trackY := y + h/2 - 3
	knobX := x + knobRadius + int(t*float32(usable)+0.5)

	win.cv.FillRoundRect(x+knobRadius, trackY, usable, 6, 3, win.theme.Border)
	win.cv.FillRoundRect(x+knobRadius, trackY, knobX-(x+knobRadius), 6, 3, win.theme.Accent)

	knob := win.theme.Accent
	if hovered || win.uiActive == id {
		knob = win.theme.AccentHover
	}
	win.cv.FillCircle(knobX, y+h/2, knobRadius, knob)
	win.cv.FillCircle(knobX, y+h/2, knobRadius-4, win.theme.Surface)
	return changed
}

// Progress draws a progress bar, with progress running from 0 to 1.
func (win *Window) Progress(x, y, w, h int, progress float32) {
	if w <= 0 || h <= 0 {
		return
	}
	progress = min(max(progress, 0), 1)
	filled := int(float32(w)*progress + 0.5)

	win.cv.FillRoundRect(x, y, w, h, h/2, win.theme.Border)
	if filled > 0 {
		win.cv.FillRoundRect(x, y, filled, h, h/2, win.theme.Accent)
	}
}

// ---------------------------------------------------------------------------
// Text field
// ---------------------------------------------------------------------------

// utf8Prev and utf8Next step the cursor by whole characters.
func utf8Prev(text string, index int) int {
	if index <= 0 {
		return 0
	}
	index--
	for index > 0 && text[index]&0xC0 == 0x80 {
		index--
	}
	return index
}

func utf8Next(text string, index int) int {
	if index >= len(text) {
		return len(text)
	}
	index++
	for index < len(text) && text[index]&0xC0 == 0x80 {
		index++
	}
	return index
}

// Input draws a single-line text field and reports whether the text changed
// this frame.
func (win *Window) Input(x, y, w, h int, text *string) bool {
	if text == nil {
		return false
	}
	const padding = 8
	id := widgetID("input", x, y, w, h, "")
	hovered := win.Hovered(x, y, w, h)
	if hovered {
		win.uiHot = id
	}

	if win.MousePressed(MouseLeft) {
		switch {
		case hovered:
			if win.uiFocus != id {
				win.uiFocus = id
				win.uiCursor = len(*text)
			}
			win.uiBlink = 0
		case win.uiFocus == id:
			win.uiFocus = 0
		}
	}
	focused := win.uiFocus == id

	changed := false
	if focused {
		changed = win.editText(text)
	}

	win.cv.FillRoundRect(x, y, w, h, win.theme.Radius, win.theme.Surface)
	border := win.theme.Border
	switch {
	case focused:
		border = win.theme.Accent
	case hovered:
		border = win.theme.TextMuted
	}
	win.cv.RoundRect(x, y, w, h, win.theme.Radius, border)

	cursor := min(win.uiCursor, len(*text))
	cursorPx := canvas.TextWidth((*text)[:cursor])
	scroll := max(cursorPx-(w-2*padding), 0)

	savedClip := win.cv.Clip
	win.cv.SetClip(x+padding, y+1, w-2*padding, h-2)

	textY := y + (h-canvas.FontHeight)/2
	win.cv.Text(x+padding-scroll, textY, *text, win.theme.Text)

	if focused && int64(win.uiBlink*2)%2 == 0 {
		win.cv.FillRect(x+padding-scroll+cursorPx, textY-1, 1, canvas.FontHeight+2, win.theme.Text)
	}

	win.cv.Clip = savedClip
	return changed
}

// editText applies this frame's keys and typing to the focused field.
func (win *Window) editText(text *string) bool {
	win.uiCursor = min(max(win.uiCursor, 0), len(*text))
	changed := false

	if win.KeyPressed(KeyBackspace) && win.uiCursor > 0 {
		start := utf8Prev(*text, win.uiCursor)
		*text = (*text)[:start] + (*text)[win.uiCursor:]
		win.uiCursor = start
		changed = true
		win.uiBlink = 0
	}
	if win.KeyPressed(KeyDelete) && win.uiCursor < len(*text) {
		end := utf8Next(*text, win.uiCursor)
		*text = (*text)[:win.uiCursor] + (*text)[end:]
		changed = true
		win.uiBlink = 0
	}
	if win.KeyPressed(KeyLeft) {
		win.uiCursor = utf8Prev(*text, win.uiCursor)
		win.uiBlink = 0
	}
	if win.KeyPressed(KeyRight) {
		win.uiCursor = utf8Next(*text, win.uiCursor)
		win.uiBlink = 0
	}
	if win.KeyPressed(KeyHome) {
		win.uiCursor = 0
	}
	if win.KeyPressed(KeyEnd) {
		win.uiCursor = len(*text)
	}
	if win.KeyPressed(KeyEscape) {
		win.uiFocus = 0
	}

	if typed := win.TextInput(); typed != "" {
		var b strings.Builder
		b.Grow(len(*text) + len(typed))
		b.WriteString((*text)[:win.uiCursor])
		b.WriteString(typed)
		b.WriteString((*text)[win.uiCursor:])
		*text = b.String()
		win.uiCursor += len(typed)
		changed = true
		win.uiBlink = 0
	}
	return changed
}