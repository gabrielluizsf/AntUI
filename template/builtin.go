package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

// builtinStyle paints exactly what the ready-made widgets paint — the same
// theme colours, the same built-in face, the same proportions — so the
// Builtin template is the look AntUI already had, at no cost. It exists to be
// the plainspoken neighbour of the futuristic one: same components, same
// events, same keyboard, different neon. Everything it draws is multiplied by
// [Scale], so the look survives a window that grows to any resolution.
type builtinStyle struct{}

func (builtinStyle) Background(win *antui.Window) canvas.Color {
	return win.Theme().Background
}

func (builtinStyle) Label(win *antui.Window, x, y int, text string) {
	win.Canvas().TextScaled(x, y, text, win.Theme().Text, Scale(win))
}

func (builtinStyle) Button(win *antui.Window, s State, x, y, w, h int, label string) {
	cv := win.Canvas()
	th := win.Theme()
	u := Scale(win)
	fill := th.Accent
	switch {
	case s.Pressed:
		fill = canvas.Shade(th.Accent, -0.18)
	case s.Hovered:
		fill = th.AccentHover
	}
	cv.FillRoundRect(x, y, w, h, th.Radius, fill)

	if s.Focused && !s.Pressed {
		cv.RoundRect(x-2*u, y-2*u, w+4*u, h+4*u, th.Radius+2*u, th.Accent)
	}

	textY := y + (h-textHeight(u))/2
	if s.Pressed {
		textY += u
	}
	cv.TextScaled(x+(w-textWidth(u, label))/2, textY, label, th.TextOnAccent, u)
}

func (builtinStyle) Checkbox(win *antui.Window, s State, x, y int, label string, on bool) {
	cv := win.Canvas()
	th := win.Theme()
	u := Scale(win)
	box := 18 * u

	fill := th.Surface
	switch {
	case on:
		fill = th.Accent
	case s.Hovered:
		fill = th.SurfaceHover
	}
	border := th.Border
	if on {
		border = th.Accent
	}
	cv.FillRoundRect(x, y, box, box, 4*u, fill)
	cv.RoundRect(x, y, box, box, 4*u, border)

	if s.Focused {
		cv.RoundRect(x-2*u, y-2*u, box+4*u, box+4*u, 6*u, th.Accent)
	}

	if on {
		mark := th.TextOnAccent
		cv.Line(x+4*u, y+9*u, x+7*u, y+13*u, mark)
		cv.Line(x+5*u, y+9*u, x+8*u, y+13*u, mark)
		cv.Line(x+7*u, y+13*u, x+14*u, y+5*u, mark)
		cv.Line(x+8*u, y+13*u, x+15*u, y+5*u, mark)
	}
	cv.TextScaled(x+box+8*u, y+(box-textHeight(u))/2, label, th.Text, u)
}

func (builtinStyle) Radio(win *antui.Window, s State, x, y int, label string, on bool) {
	cv := win.Canvas()
	th := win.Theme()
	u := Scale(win)
	size := 18 * u

	surface := th.Surface
	if s.Hovered {
		surface = th.SurfaceHover
	}
	border := th.Border
	if on {
		border = th.Accent
	}
	cv.FillCircle(x+size/2, y+size/2, size/2, surface)
	cv.Circle(x+size/2, y+size/2, size/2, border)
	if s.Focused {
		cv.Circle(x+size/2, y+size/2, size/2+3*u, th.Accent)
	}
	if on {
		cv.FillCircle(x+size/2, y+size/2, size/2-4*u, th.Accent)
	}
	cv.TextScaled(x+size+8*u, y+(size-textHeight(u))/2, label, th.Text, u)
}

func (builtinStyle) Slider(win *antui.Window, s State, x, y, w, h int, value float32) {
	cv := win.Canvas()
	th := win.Theme()
	u := Scale(win)

	value = min(max(value, 0), 1)

	// Keep the thumb completely inside the control.
	knobRadius := max(h/2, 7*u)

	if w <= knobRadius*2 {
		return
	}

	cy := y + h/2

	// The track goes from the center of the left thumb position
	// to the center of the right thumb position.
	trackX := x + knobRadius
	trackW := max(w-2*knobRadius, 1)

	trackH := 8 * u
	trackY := cy - trackH/2

	knobX := trackX + int(value*float32(trackW)+0.5)
	knobX = min(max(knobX, trackX), trackX+trackW)

	// ---------------------------------------------------------------------
	// Track
	// ---------------------------------------------------------------------

	// Unfilled part.
	cv.FillRoundRect(
		trackX,
		trackY,
		trackW,
		trackH,
		trackH/2,
		th.Border,
	)

	// Filled part.
	fillW := knobX - trackX
	if fillW > 0 {
		cv.FillRoundRect(
			trackX,
			trackY,
			fillW,
			trackH,
			trackH/2,
			th.Accent,
		)
	}

	// ---------------------------------------------------------------------
	// Thumb
	// ---------------------------------------------------------------------

	knob := th.Accent

	if s.Active {
		knob = th.AccentHover
	} else if s.Hovered {
		knob = th.AccentHover
	}

	// Outer thumb.
	cv.FillCircle(
		knobX,
		cy,
		knobRadius,
		knob,
	)

	// Inner surface creates the visual "thumb" ring.
	innerRadius := max(knobRadius-4*u, 2*u)

	cv.FillCircle(
		knobX,
		cy,
		innerRadius,
		th.Surface,
	)

	// Focus ring.
	if s.Focused {
		cv.Circle(
			knobX,
			cy,
			knobRadius+3*u,
			th.Accent,
		)
	}
}

func (builtinStyle) Input(win *antui.Window, s State, x, y, w, h int, text string, cursor int, caret bool) {
	cv := win.Canvas()
	th := win.Theme()
	u := Scale(win)
	padding := 8 * u

	cv.FillRoundRect(x, y, w, h, th.Radius, th.Surface)
	border := th.Border
	switch {
	case s.Focused:
		border = th.Accent
	case s.Hovered:
		border = th.TextMuted
	}
	cv.RoundRect(x, y, w, h, th.Radius, border)

	cursor = min(max(cursor, 0), len(text))
	cursorPx := canvas.TextWidth(text[:cursor]) * u
	scroll := max(cursorPx-(w-2*padding), 0)

	savedClip := cv.Clip
	cv.SetClip(x+padding, y+u, w-2*padding, h-2*u)

	textY := y + (h-textHeight(u))/2
	cv.TextScaled(x+padding-scroll, textY, text, th.Text, u)

	if s.Focused && caret {
		cv.FillRect(x+padding-scroll+cursorPx, textY-u, u, textHeight(u)+2*u, th.Text)
	}
	cv.Clip = savedClip
}

func (builtinStyle) Select(win *antui.Window, s State, x, y, w, h int, value string, open bool) {
	cv := win.Canvas()
	th := win.Theme()
	u := Scale(win)

	local := th.Surface
	if s.Hovered {
		local = th.SurfaceHover
	}
	border := th.Border
	if s.Focused || open {
		border = th.Accent
	}
	cv.FillRoundRect(x, y, w, h, th.Radius, local)
	cv.RoundRect(x, y, w, h, th.Radius, border)

	padding := 8 * u
	textY := y + (h-textHeight(u))/2
	cv.TextScaled(x+padding, textY, value, th.Text, u)

	// The arrow that says "there is more below".
	aw, ah := 6*u, 4*u
	ax := x + w - padding - aw
	ay := y + (h-ah)/2
	cv.Line(ax, ay, ax+aw/2, ay+ah, th.TextMuted)
	cv.Line(ax+aw/2, ay+ah, ax+aw, ay, th.TextMuted)
}

func (builtinStyle) SelectOption(win *antui.Window, s State, x, y, w, h int, label string, selected bool) {
	cv := win.Canvas()
	th := win.Theme()

	fill := th.Surface
	border := th.Border
	if selected {
		fill = th.SurfaceHover
		border = th.Accent
	}
	if s.Hovered {
		fill = th.AccentHover
		border = th.Accent
	}
	cv.FillRoundRect(x, y, w, h, 0, fill)
	if selected || s.Hovered {
		cv.FillRect(x, y, 3, h, th.Accent)
	}
	cv.RoundRect(x, y, w, h, 0, border)

	u := Scale(win)
	textY := y + (h-textHeight(u))/2
	colour := th.Text
	if selected {
		colour = th.TextOnAccent
	}
	cv.TextScaled(x+10*u, textY, label, colour, u)
}

func (builtinStyle) TextArea(win *antui.Window, s State, x, y, w, h int, text string, cursor int, caret bool) {
	cv := win.Canvas()
	th := win.Theme()
	u := Scale(win)
	padding := 8 * u

	cv.FillRoundRect(x, y, w, h, th.Radius, th.Surface)
	border := th.Border
	switch {
	case s.Focused:
		border = th.Accent
	case s.Hovered:
		border = th.TextMuted
	}
	cv.RoundRect(x, y, w, h, th.Radius, border)

	innerW := w - 2*padding
	textW := max(innerW/u, 12)
	cursor = min(max(cursor, 0), len(text))
	lines := textLines(text, textW)

	lineH := textHeight(u)
	visible := max((h-2*padding)/lineH, 1)
	viewTop := 0
	for li, l := range lines {
		if cursor >= l[0] && cursor < l[1] || (li == len(lines)-1 && cursor == l[1]) {
			viewTop = max(li-visible+1, 0)
			break
		}
	}
	viewTop = max(min(viewTop, len(lines)-visible), 0)

	savedClip := cv.Clip
	cv.SetClip(x+padding, y+padding, innerW, h-2*padding)

	for lif, l := range lines {
		if lif < viewTop || lif >= viewTop+visible {
			continue
		}
		cy := y + padding + (lif-viewTop)*lineH
		line := text[l[0]:l[1]]
		cv.TextScaled(x+padding, cy, line, th.Text, u)
		if s.Focused && caret && cursor >= l[0] && cursor <= l[1] {
			off := cursor - l[0]
			cx := x + padding + textWidth(u, line[:off])
			cv.FillRect(cx, cy, u, lineH, th.Text)
		}
	}
	cv.Clip = savedClip
}

func (builtinStyle) InputTextPos(win *antui.Window, s State, x, y, w, h int, text string, cursor, mx, my int) int {
	u := Scale(win)
	padding := 8 * u
	cursor = min(max(cursor, 0), len(text))
	cursorPx := textWidth(u, text[:cursor])
	scroll := max(cursorPx-(w-2*padding), 0)
	return caretAt(func(prefix string) int { return textWidth(u, prefix) }, text, x+padding-scroll, mx)
}

func (builtinStyle) TextAreaTextPos(win *antui.Window, s State, x, y, w, h int, text string, cursor, mx, my int) int {
	u := Scale(win)
	padding := 8 * u
	down := max((w-2*padding)/u, 12)
	cursor = min(max(cursor, 0), len(text))
	lines := textLines(text, down)
	return textAreaTextPos(text, lines, padding, x, y, h, cursor, textHeight(u),
		func(prefix string) int { return textWidth(u, prefix) }, mx, my)
}

func (builtinStyle) Switch(win *antui.Window, s State, x, y int, label string, on bool) {
	cv := win.Canvas()
	th := win.Theme()
	u := Scale(win)
	trackW, trackH := 36*u, 20*u

	fill := th.Surface
	if s.Hovered {
		fill = th.SurfaceHover
	}
	border := th.Border
	if on || s.Focused {
		border = th.Accent
	}
	cv.FillRoundRect(x, y, trackW, trackH, trackH/2, fill)
	cv.RoundRect(x, y, trackW, trackH, trackH/2, border)

	colour := th.Border
	if on {
		colour = th.Accent
	}
	cx := x + trackH/2
	if on {
		cx = x + trackW - trackH/2
	}
	cv.FillCircle(cx, y+trackH/2, trackH/2-2*u, colour)
	if s.Focused {
		cv.Circle(x+trackW/2, y+trackH/2, trackH/2+3*u, th.Accent)
	}
	cv.TextScaled(x+trackW+8*u, y+(trackH-textHeight(u))/2, label, th.Text, u)
}

func (builtinStyle) Progress(win *antui.Window, x, y, w, h int, value float32) {
	cv := win.Canvas()
	th := win.Theme()
	u := Scale(win)

	value = min(max(value, 0), 1)
	filled := int(float32(w)*value + 0.5)

	cv.FillRoundRect(x, y, w, h, h/2, th.Border)
	if filled > 0 {
		cv.FillRoundRect(x, y, filled, h, h/2, th.Accent)
	}
	_ = u
}

func (builtinStyle) DatePicker(win *antui.Window, s State, x, y, w, h int, year, month, firstWD, days, selected, today, hover int) {
	cv := win.Canvas()
	th := win.Theme()
	u := Scale(win)

	cv.FillRoundRect(x, y, w, h, 4*u, th.Surface)
	border := th.Border
	if s.Focused {
		border = th.Accent
	}
	cv.RoundRect(x, y, w, h, 4*u, border)

	// Header: the prev/next arrows and the month name.
	mname := monthName(month)
	header := dateHeaderHeight(u)
	aw := 8 * u
	cv.TextScaled(x+w/2-textWidth(u, mname)/2, y+(header-textHeight(u))/2, mname, th.Text, u)
	cv.TextScaled(x+w/2-textWidth(u, mname)/2+textWidth(u, mname)+4*u, y+(header-textHeight(u))/2,
		itoa(year), th.TextMuted, u)

	// Arrow glyphs.
	prevX, nextX := x+2*u, x+w-2*u-aw
	ay := y + (header-textHeight(u))/2
	cv.Line(prevX+3*u, ay+textHeight(u)/2, prevX+aw-3*u, ay+3*u, th.TextMuted)
	cv.Line(prevX+3*u, ay+textHeight(u)/2, prevX+aw-3*u, ay+textHeight(u)-3*u, th.TextMuted)
	cv.Line(nextX+3*u, ay+3*u, nextX+aw-3*u, ay+textHeight(u)/2, th.TextMuted)
	cv.Line(nextX+aw-3*u, ay+textHeight(u)/2, nextX+3*u, ay+textHeight(u)-3*u, th.TextMuted)

	// Column heads S M T W T F S.
	c := dateCellSize(u)
	gx, _, _, _ := dateGridRect(x, y, u)
	hy := weekdayHeadY(x, y, u)
	for d := 0; d < 7; d++ {
		cx := gx + d*c
		cv.TextScaled(cx+(c-textWidth(u, "S"))/2, hy, dayHeadName((firstWD+d)%7), th.TextMuted, u)
	}

	// Day cells.
	for day := 1; day <= days; day++ {
		pos := firstWD + day - 1
		col, row := pos%7, pos/7
		cx, cy, cc, _ := dayCellRect(x, y, u, col, row)
		switch day {
		case selected:
			cv.FillRoundRect(cx, cy, cc, cc, 2*u, th.Accent)
			cv.TextScaled(cx+(cc-textWidth(u, itoa(day)))/2, cy+(cc-textHeight(u))/2, itoa(day), th.TextOnAccent, u)
		case hover:
			cv.FillRoundRect(cx, cy, cc, cc, 2*u, th.SurfaceHover)
			cv.TextScaled(cx+(cc-textWidth(u, itoa(day)))/2, cy+(cc-textHeight(u))/2, itoa(day), th.Text, u)
		default:
			colour := th.Text
			if day == today {
				colour = th.Accent
			}
			cv.TextScaled(cx+(cc-textWidth(u, itoa(day)))/2, cy+(cc-textHeight(u))/2, itoa(day), colour, u)
		}
		if day == today && day != selected {
			cv.Circle(cx+cc/2, cy+cc/2, cc/2, th.Accent)
		}
	}
}

func (builtinStyle) DatePickerBox(win *antui.Window, s State, x, y, w, h int, value string, open bool) {
	cv := win.Canvas()
	th := win.Theme()
	u := Scale(win)

	local := th.Surface
	if s.Hovered {
		local = th.SurfaceHover
	}
	border := th.Border
	if s.Focused || open {
		border = th.Accent
	}
	cv.FillRoundRect(x, y, w, h, th.Radius, local)
	cv.RoundRect(x, y, w, h, th.Radius, border)

	padding := 8 * u
	cv.TextScaled(x+padding, y+(h-textHeight(u))/2, value, th.Text, u)

	// The caret that says "there is a calendar below".
	aw, ah := 6*u, 4*u
	ax := x + w - padding - aw
	ay := y + (h-ah)/2
	if open {
		cv.Line(ax, ay+ah, ax+aw/2, ay, th.TextMuted)
		cv.Line(ax+aw/2, ay, ax+aw, ay+ah, th.TextMuted)
	} else {
		cv.Line(ax, ay, ax+aw/2, ay+ah, th.TextMuted)
		cv.Line(ax+aw/2, ay+ah, ax+aw, ay, th.TextMuted)
	}
}

func monthName(m int) string {
	switch m {
	case 1:
		return "January"
	case 2:
		return "February"
	case 3:
		return "March"
	case 4:
		return "April"
	case 5:
		return "May"
	case 6:
		return "June"
	case 7:
		return "July"
	case 8:
		return "August"
	case 9:
		return "September"
	case 10:
		return "October"
	case 11:
		return "November"
	}
	return "December"
}

func dayHeadName(wd int) string {
	return "SMTWTFS"[wd : wd+1]
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [12]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
