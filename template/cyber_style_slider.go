package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

// Slider paints a precision control: a thin lit trough, graduations along
// it, and a knob that is a ring around a bright core — a reticle, the same
// language the rest of the theme speaks. The filled run gives off a faint
// glow of its own, so it reads as powered rather than painted.
func (cyberStyle) Slider(win *antui.Window, s State, x, y, w, h int, value float32) {
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)

	value = min(max(value, 0), 1)

	// Keep the knob completely inside the slider.
	knobRadius := max(h/2, 8*u)

	if w <= knobRadius*2 {
		return
	}

	cy := y + h/2

	// Track runs between the centers of the knob at both extremes.
	trackX := x + knobRadius
	trackW := max(w-2*knobRadius, 1)

	trackH := 8 * u
	trackY := cy - trackH/2

	knobX := trackX + int(value*float32(trackW)+0.5)
	knobX = min(max(knobX, trackX), trackX+trackW)

	fillW := knobX - trackX

	// ---------------------------------------------------------------------
	// Track glow
	// ---------------------------------------------------------------------

	// Very subtle glow behind the entire track.
	glow(
		cv,
		trackX-2*u,
		trackY-2*u,
		trackW+4*u,
		trackH+4*u,
		trackH/2+2*u,
		v.neon,
		5,
		u,
	)

	// Dark glass body.
	cv.FillRoundRect(
		trackX,
		trackY,
		trackW,
		trackH,
		trackH/2,
		canvas.Fade(cyberPalette.panelAct, 230),
	)

	// Thin outer edge.
	cv.RoundRect(
		trackX,
		trackY,
		trackW,
		trackH,
		trackH/2,
		v.edge,
	)

	// ---------------------------------------------------------------------
	// Filled track
	// ---------------------------------------------------------------------

	if fillW > 0 {
		// Glow around the active portion.
		glow(
			cv,
			trackX-u,
			trackY-u,
			fillW+2*u,
			trackH+2*u,
			trackH/2+u,
			v.neon,
			10,
			u,
		)

		cv.FillRoundRect(
			trackX,
			trackY,
			fillW,
			trackH,
			trackH/2,
			v.neon,
		)
	}

	// ---------------------------------------------------------------------
	// Scale marks
	// ---------------------------------------------------------------------

	for i := 1; i < 4; i++ {
		gx := trackX + trackW*i/4

		tickColor := canvas.Fade(v.dim, 130)

		if gx <= knobX {
			tickColor = v.neon
		}

		cv.FillRoundRect(
			gx-u,
			trackY-3*u,
			2*u,
			trackH+6*u,
			u,
			tickColor,
		)
	}

	// ---------------------------------------------------------------------
	// End terminals
	// ---------------------------------------------------------------------

	for _, px := range []int{
		trackX,
		trackX + trackW,
	} {
		cv.FillCircle(
			px,
			cy,
			3*u,
			cyberPalette.panelAct,
		)

		cv.Circle(
			px,
			cy,
			3*u,
			v.edge,
		)

		cv.FillCircle(
			px,
			cy,
			u,
			v.neon,
		)
	}

	// ---------------------------------------------------------------------
	// Thumb
	// ---------------------------------------------------------------------

	halo := 22

	if s.Hovered {
		halo = 40
	}

	if s.Active {
		halo = 70
	}

	// Large but soft glow around the knob.
	glow(
		cv,
		knobX-knobRadius-3*u,
		cy-knobRadius-3*u,
		knobRadius*2+6*u,
		knobRadius*2+6*u,
		knobRadius+2*u,
		v.neon,
		int(float64(halo)*v.glow),
		u,
	)

	// Outer neon shell.
	cv.FillCircle(
		knobX,
		cy,
		knobRadius+2*u,
		v.neon,
	)

	// Dark body.
	cv.FillCircle(
		knobX,
		cy,
		knobRadius,
		cyberPalette.panelAct,
	)

	// Bright inner ring.
	cv.Circle(
		knobX,
		cy,
		knobRadius-2*u,
		v.neon,
	)

	// Inner core.
	cv.FillCircle(
		knobX,
		cy,
		3*u,
		v.neon,
	)

	// Focus gives the thumb an additional outer ring.
	if s.Focused {
		cv.Circle(
			knobX,
			cy,
			knobRadius+5*u,
			v.text,
		)
	}
}
