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
