package template

import "github.com/gabrielluizsf/antui"

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
