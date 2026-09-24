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

func (cyberStyle) Select(win *antui.Window, s State, x, y, w, h int, value string, open bool) {
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)
	radius := buttonRadius(h, u)

	body, border := cyberPalette.panel, v.edge
	halo := 0
	if s.Focused || open {
		body, border, halo = cyberPalette.panel, v.neon, 18
	}
	if s.Hovered {
		body = cyberPalette.panelHov
	}
	if halo > 0 {
		glow(cv, x, y, w, h, radius, v.neon, int(float64(halo)*v.glow), u)
	}
	neonFill(cv, x, y, w, h, radius, body, border)

	padding := 8 * u
	textY := y + (h-textHeight(u))/2
	colour := v.text
	if value == "" {
		colour = v.dim
	}
	cv.TextScaled(x+padding, textY, value, colour, u)

	// The caret arrow.
	aw, ah := 6*u, 4*u
	ax := x + w - padding - aw
	ay := y + (h-ah)/2
	if open {
		cv.Line(ax, ay+ah, ax+aw/2, ay, v.neon)
		cv.Line(ax+aw/2, ay, ax+aw, ay+ah, v.neon)
	} else {
		cv.Line(ax, ay, ax+aw/2, ay+ah, v.dim)
		cv.Line(ax+aw/2, ay+ah, ax+aw, ay, v.dim)
	}
}

func (cyberStyle) SelectOption(win *antui.Window, s State, x, y, w, h int, label string, selected bool) {
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)

	body := cyberPalette.panel
	if selected || s.Hovered {
		body = cyberPalette.panelHov
	}
	cv.FillRect(x, y, w, h, body)
	edge := v.edge
	if selected {
		edge = v.neon
	}
	cv.Line(x, y, x+w, y, edge)
	if selected || s.Hovered {
		cv.FillRect(x, y, 3*u, h, v.neon)
	}
	colour := v.text
	if selected {
		colour = v.neon
	}
	cv.TextScaled(x+10*u, y+(h-textHeight(u))/2, label, colour, u)
}

func (cyberStyle) TextArea(win *antui.Window, s State, x, y, w, h int, text string, cursor int, caret bool) {
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
		colour := v.text
		if line == "" {
			colour = v.dim
		}
		cv.TextScaled(x+padding, cy, line, colour, u)
		if s.Focused && caret && cursor >= l[0] && cursor <= l[1] {
			off := cursor - l[0]
			cx := x + padding + textWidth(u, line[:off])
			cv.FillRect(cx, cy, 2*u, lineH, v.neon)
			cv.FillRect(cx, cy+lineH, 2*u, 2*u, canvas.Fade(v.neon, 150))
		}
	}
	cv.Clip = savedClip
}

func (cyberStyle) InputTextPos(win *antui.Window, s State, x, y, w, h int, text string, cursor, mx, my int) int {
	u := Scale(win)
	padding := 8 * u
	cursor = min(max(cursor, 0), len(text))
	cursorPx := textWidth(u, text[:cursor])
	scroll := max(cursorPx-(w-2*padding), 0)
	return caretAt(func(prefix string) int { return textWidth(u, prefix) }, text, x+padding-scroll, mx)
}

func (cyberStyle) TextAreaTextPos(win *antui.Window, s State, x, y, w, h int, text string, cursor, mx, my int) int {
	u := Scale(win)
	padding := 8 * u
	cursor = min(max(cursor, 0), len(text))
	lines := textLines(text, max((w-2*padding)/u, 12))
	return textAreaTextPos(text, lines, padding, x, y, h, cursor, textHeight(u),
		func(prefix string) int { return textWidth(u, prefix) }, mx, my)
}

func (cyberStyle) Switch(win *antui.Window, s State, x, y int, label string, on bool) {
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)
	trackW, trackH := 36*u, 20*u

	if s.Focused || s.Hovered {
		glow(cv, x, y, trackW, trackH, trackH, v.neon, int(float64(30)*v.glow), u)
	}
	body := cyberPalette.panel
	if s.Hovered {
		body = cyberPalette.panelHov
	}
	border := v.edge
	if on || s.Focused {
		border = v.neon
	}
	cv.FillRoundRect(x, y, trackW, trackH, trackH/2, body)
	cv.RoundRect(x, y, trackW, trackH, trackH/2, border)

	if on {
		cv.FillRoundRect(x, y, trackH, trackH, trackH/2, canvas.Fade(v.neon, 120))
	}
	colour := v.edge
	if on {
		colour = v.neon
	}
	cx := x + trackH/2
	if on {
		cx = x + trackW - trackH/2
	}
	cv.FillCircle(cx, y+trackH/2, trackH/2-3*u, colour)
	cv.Circle(cx, y+trackH/2, trackH/2-4*u, cyberPalette.panelAct)
	cv.TextScaled(x+trackW+8*u, y+(trackH-textHeight(u))/2, label, v.text, u)
}

func (cyberStyle) Progress(win *antui.Window, x, y, w, h int, value float32) {
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)

	value = min(max(value, 0), 1)
	filled := int(float32(w)*value + 0.5)

	cv.FillRoundRect(x, y, w, h, h/2, cyberPalette.panelAct)
	cv.RoundRect(x, y, w, h, h/2, v.edge)
	if filled > 0 {
		glow(cv, x, y, filled, h, h/2, v.neon, 5, u)
		cv.FillRoundRect(x, y, filled, h, h/2, v.neon)
	}
}

func (cyberStyle) DatePicker(win *antui.Window, s State, x, y, w, h int, year, month, firstWD, days, selected, today, hover int) {
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)

	glow(cv, x, y, w, h, 4*u, v.neon, 12, u)
	body := cyberPalette.panel
	if s.Hovered {
		body = cyberPalette.panelHov
	}
	neonFill(cv, x, y, w, h, 4*u, body, v.edge)

	// Header: month and year, flanked by the arrows.
	header := dateHeaderHeight(u)
	mname := monthName(month)
	labelW := textWidth(u, mname) + 3*u + textWidth(u, itoa(year))
	lx := x + (w-labelW)/2
	cv.TextScaled(lx, y+(header-textHeight(u))/2, mname, v.text, u)
	cv.TextScaled(lx+textWidth(u, mname)+3*u, y+(header-textHeight(u))/2, itoa(year), v.dim, u)
	cv.Line(x+3*u, y+header/2, x+6*u, y+header/2, v.neon)
	cv.Line(x+2*u, y+header/2, x+5*u, y+header/2, v.neon)
	cv.Line(x+w-3*u, y+header/2, x+w-6*u, y+header/2, v.neon)
	cv.Line(x+w-2*u, y+header/2, x+w-5*u, y+header/2, v.neon)

	// Column heads.
	c := dateCellSize(u)
	gx, _, _, _ := dateGridRect(x, y, u)
	hy := weekdayHeadY(x, y, u)
	for d := 0; d < 7; d++ {
		cx := gx + d*c
		cv.TextScaled(cx+(c-textWidth(u, "S"))/2, hy, dayHeadName((firstWD+d)%7), v.dim, u)
	}

	// Day cells.
	for day := 1; day <= days; day++ {
		pos := firstWD + day - 1
		col, row := pos%7, pos/7
		cx, cy, cc, _ := dayCellRect(x, y, u, col, row)
		switch day {
		case selected:
			glow(cv, cx, cy, cc, cc, 2*u, v.neon, int(float64(20)*v.glow), u)
			cv.FillRoundRect(cx, cy, cc, cc, 2*u, v.neon)
			cv.TextScaled(cx+(cc-textWidth(u, itoa(day)))/2, cy+(cc-textHeight(u))/2, itoa(day), cyberPalette.bg, u)
		case hover:
			cv.FillRoundRect(cx, cy, cc, cc, 2*u, cyberPalette.panelHov)
			cv.TextScaled(cx+(cc-textWidth(u, itoa(day)))/2, cy+(cc-textHeight(u))/2, itoa(day), v.text, u)
		default:
			colour := v.text
			if day == today {
				colour = v.neon
			}
			cv.TextScaled(cx+(cc-textWidth(u, itoa(day)))/2, cy+(cc-textHeight(u))/2, itoa(day), colour, u)
		}
		if day == today && day != selected {
			cv.Circle(cx+cc/2, cy+cc/2, cc/2-2*u, v.neon)
		}
	}
}

func (cyberStyle) DatePickerBox(win *antui.Window, s State, x, y, w, h int, value string, open bool) {
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)
	radius := buttonRadius(h, u)

	body, border := cyberPalette.panel, v.edge
	halo := 0
	if s.Focused || open {
		body, border, halo = cyberPalette.panel, v.neon, 18
	}
	if s.Hovered {
		body = cyberPalette.panelHov
	}
	if halo > 0 {
		glow(cv, x, y, w, h, radius, v.neon, int(float64(halo)*v.glow), u)
	}
	neonFill(cv, x, y, w, h, radius, body, border)

	padding := 8 * u
	colour := v.text
	if value == "" {
		colour = v.dim
	}
	cv.TextScaled(x+padding, y+(h-textHeight(u))/2, value, colour, u)

	// The caret arrow.
	aw, ah := 6*u, 4*u
	ax := x + w - padding - aw
	ay := y + (h-ah)/2
	if open {
		cv.Line(ax, ay+ah, ax+aw/2, ay, v.neon)
		cv.Line(ax+aw/2, ay, ax+aw, ay+ah, v.neon)
	} else {
		cv.Line(ax, ay, ax+aw/2, ay+ah, v.dim)
		cv.Line(ax+aw/2, ay+ah, ax+aw, ay, v.dim)
	}
}
