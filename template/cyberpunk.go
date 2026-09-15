package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

// cyberStyle paints the futuristic look: dark glass panels, thin neon edges,
// and a soft glow around whatever the pointer is over. There are no gradients
// on the canvas, so the glow is faked the honest way — concentric shells of
// translucent colour, the same way the glow sits on real neon. Like every
// style, everything it draws is multiplied by [Scale], so the design holds at
// any resolution.
type cyberStyle struct{}

// cyberPalette is the whole range of the futuristic look.
var cyberPalette = struct {
	bg       canvas.Color
	panel    canvas.Color
	panelHov canvas.Color
	panelAct canvas.Color
	edge     canvas.Color
	neon     canvas.Color
	text     canvas.Color
	textDim  canvas.Color
}{
	bg:       0xFF0B0E1A,
	panel:    0xFF12182B,
	panelHov: 0xFF17213D,
	panelAct: 0xFF0D1222,
	edge:     0xFF24365C,
	neon:     0xFF3DE8FF,
	text:     0xFFE8F0FF,
	textDim:  0xFF8FA2C0,
}

func (cyberStyle) Background(win *antui.Window) canvas.Color {
	return cyberPalette.bg
}

// vivid is what the theme draws the ink in, sharpened for the window's
// opacity: text and the neon details get lifted toward white as the window
// fades, so a translucent window — whose whole frame the compositor dims —
// keeps its fonts and its button details bright and readable. The dark panels
// are left alone: the contrast between ink and panel is exactly what the lift
// is saving. At full opacity it is the palette as it sits.
type vivid struct {
	text, neon, dim, edge canvas.Color
	glow                  float64 // how much the halos are amplified, ≥ 1
}

// vividFor reads the window's opacity and works out how loudly to draw. The
// boost is 255/opacity, one at full opacity and three at a third — the
// brightest ink a colour can survive.
func vividFor(win *antui.Window) vivid {
	o := int(win.Opacity())
	f := 1.0
	if o > 0 {
		f = 255.0 / float64(o)
	}
	if f > 3 {
		f = 3
	}
	return vivid{
		text: lift(cyberPalette.text, f),
		neon: lift(cyberPalette.neon, f),
		dim:  lift(cyberPalette.textDim, f),
		edge: lift(cyberPalette.edge, f),
		glow: f,
	}
}

// lift brings a colour partway to white. The weight is (f-1)/f of the way:
// identity at f = 1, halfway at 50% opacity, capped two thirds of the way
// there. RGB climbs, alpha stays.
func lift(c canvas.Color, f float64) canvas.Color {
	if f <= 1 {
		return c
	}
	alpha := byte(c >> 24)
	w := int((f - 1) * 255 / f)
	up := func(v byte) byte {
		return byte(int(v) + (255-int(v))*w/255)
	}
	return canvas.Color(alpha)<<24 |
		canvas.Color(up(byte(c>>16)))<<16 |
		canvas.Color(up(byte(c>>8)))<<8 |
		canvas.Color(up(byte(c)))
}

func (cyberStyle) Label(win *antui.Window, x, y int, text string) {
	win.Canvas().TextScaled(x, y, text, vividFor(win).text, Scale(win))
}

// glow draws translucent shells around a rounded rectangle, the strongest
// closest to it, spreading colour outward like light. u is the scale, so the
// shells stay proportioned to the control on every resolution.
func glow(cv *canvas.Canvas, x, y, w, h, radius int, c canvas.Color, strength, u int) {
	for i := 2; i >= 0; i-- {
		alpha := strength / (3 - i)
		if alpha <= 0 {
			continue
		}
		g := (i + 1) * u
		cv.FillRoundRect(x-g, y-g, w+g*2, h+g*2, radius+g, canvas.Fade(c, alpha))
	}
}

// neonFill paints a control's body and its border: a dark panel with a thin
// coloured edge. What happens under the pointer — who glows, who brightens —
// is the caller's call, because a button and a field answer differently.
func neonFill(cv *canvas.Canvas, x, y, w, h, radius int, body, border canvas.Color) {
	cv.FillRoundRect(x, y, w, h, radius, body)
	cv.RoundRect(x, y, w, h, radius, border)
}

// buttonRadius is how rounded a control's corners are. It is capped at eight
// scale units, so a button never turns into a pill: just enough rounding to
// sit like the modern dashboard, whatever its height.
func buttonRadius(h, u int) int { return min(h/2, 8*u) }

func (cyberStyle) Button(win *antui.Window, s State, x, y, w, h int, label string) {
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)
	radius := buttonRadius(h, u)

	// Idle is a quiet panel with a blue edge. Attention lights it: the border
	// turns neon and a soft halo grows around it, larger while it is held.
	body, border := cyberPalette.panel, v.edge
	halo := 0
	switch {
	case s.Pressed:
		body, border, halo = cyberPalette.panelAct, v.neon, 30
	case s.Hovered:
		body, border, halo = cyberPalette.panelHov, v.neon, 18
	case s.Focused:
		border, halo = v.neon, 16
	}
	if halo > 0 {
		glow(cv, x, y, w, h, radius, v.neon, int(float64(halo)*v.glow), u)
	}
	neonFill(cv, x, y, w, h, radius, body, border)
	cv.TextScaled(x+(w-textWidth(u, label))/2, y+(h-textHeight(u))/2, label, v.text, u)
}

func (cyberStyle) Checkbox(win *antui.Window, s State, x, y int, label string, on bool) {
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)
	box := 18 * u
	radius := 4 * u
	if s.Focused || s.Hovered {
		glow(cv, x, y, box, box, radius, v.neon, int(float64(30)*v.glow), u)
	}
	body := cyberPalette.panel
	if s.Hovered {
		body = cyberPalette.panelHov
	}
	cv.FillRoundRect(x, y, box, box, radius, body)
	border := v.edge
	if on || s.Hovered {
		border = v.neon
	}
	cv.RoundRect(x, y, box, box, radius, border)
	if on {
		cv.FillRect(x+4*u, y+9*u, 3*u, 3*u, v.neon)
		cv.FillRect(x+6*u, y+9*u, 8*u, 3*u, v.neon)
	}
	cv.TextScaled(x+box+8*u, y+(box-textHeight(u))/2, label, v.text, u)
}

func (cyberStyle) Radio(win *antui.Window, s State, x, y int, label string, on bool) {
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)
	size := 18 * u
	cx, cy := x+size/2, y+size/2
	if s.Focused || s.Hovered {
		glow(cv, x, y, size, size, size/2, v.neon, int(float64(30)*v.glow), u)
	}
	body := cyberPalette.panel
	if s.Hovered {
		body = cyberPalette.panelHov
	}
	cv.FillCircle(cx, cy, size/2, body)
	border := v.edge
	if on || s.Hovered {
		border = v.neon
	}
	cv.Circle(cx, cy, size/2, border)
	if on {
		cv.FillCircle(cx, cy, size/2-4*u, v.neon)
	}
	cv.TextScaled(x+size+8*u, y+(size-textHeight(u))/2, label, v.text, u)
}

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

func (cyberStyle) Input(win *antui.Window, s State, x, y, w, h int, text string, cursor int, caret bool) {
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)
	radius := buttonRadius(h, u)

	body, border := cyberPalette.panel, v.edge
	if s.Focused {
		border = v.neon
		glow(cv, x, y, w, h, radius, v.neon, int(float64(20)*v.glow), u)
	} else if s.Hovered {
		body, border = cyberPalette.panelHov, v.neon
	}
	neonFill(cv, x, y, w, h, radius, body, border)

	padding := 8 * u
	cursor = min(max(cursor, 0), len(text))
	cursorPx := canvas.TextWidth(text[:cursor]) * u
	scroll := max(cursorPx-(w-2*padding), 0)

	savedClip := cv.Clip
	cv.SetClip(x+padding, y+u, w-2*padding, h-2*u)

	textY := y + (h-textHeight(u))/2
	colour := v.text
	if text == "" {
		colour = v.dim
	}
	cv.TextScaled(x+padding-scroll, textY, text, colour, u)
	if s.Focused && caret {
		cx := x + padding - scroll + cursorPx
		cv.FillRect(cx, textY, 2*u, textHeight(u), v.neon)
		cv.FillRect(cx, textY+textHeight(u), 2*u, 2*u, canvas.Fade(v.neon, 150))
	}
	cv.Clip = savedClip
}
