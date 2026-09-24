package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// cssStyle paints the CSS template: every seam is a decision from a computed
// [css.Style], and everything the stylesheet did not say falls back onto the
// window's theme. It is the [Style] behind [TemplateWithCSS].
type cssStyle struct {
	win     *antui.Window
	classes *css.CSSClasses
	images  map[string]*canvas.Canvas // background-image url() targets

	// The widget timeline: one entry per widget as it appears in a frame, so
	// a widget's state changes and their transitions persist across Reset
	// calls. Entries are keyed by role, ordinal in the frame and label;
	// entryStyle blends them while they animate.
	timeline map[string]*widgetEntry
	cur      *widgetEntry   // the widget laid out right now
	seq      int            // widgets begun so far in this frame
	frame    int            // Reset calls carried out so far
	now      func() float64 // clock in seconds, antui.Now by default

	// transformed widgets first paint to a window-sized scratch layer, then
	// composite through their matrix; considered dirty every frame.
	layer  *canvas.Canvas
	layerW int
	layerH int
}

// newCSSStyle makes a cssStyle over an empty class table; the table is shared
// with its template, so SetStyle on the template reaches the painter.
func newCSSStyle(win *antui.Window) *cssStyle {
	return &cssStyle{
		win:      win,
		classes:  css.NewTable(),
		images:   map[string]*canvas.Canvas{},
		timeline: map[string]*widgetEntry{},
		now:      antui.Now,
	}
}

// image is one registered background-image. A name the stylesheet asks for
// that was never registered paints nothing.
func (cs *cssStyle) image(name string) *canvas.Canvas { return cs.images[name] }

// styleFor resolves the animated style for a widget role, folding the
// interaction state into the pseudo-classes. While a widget is being laid
// out, styleFor drives its timeline entry; otherwise the plain cascade
// answer is returned, which keeps styles the template reads for itself free
// of any frame's animation.
func (cs *cssStyle) styleFor(role string, s State) css.Style {
	if cs.cur != nil && cs.cur.role == role {
		return cs.entryStyle(cs.cur, role, s)
	}
	return cs.baseStyle(role, s)
}

// baseStyle is the computed cascade style for a role and state, animations
// not yet folded in.
func (cs *cssStyle) baseStyle(role string, s State) css.Style {
	return cs.classes.GetStyle(role, nil, css.StateWith(
		s.Hovered, s.Focused, s.Active || s.Pressed, s.On,
	), cs.win.Width())
}

// unitsCtx is the measurement context an animated style resolves lengths
// and viewport units in: the window and the widget's own font size.
func (cs *cssStyle) unitsCtx(st css.Style) css.Units {
	u := css.Units{Width: cs.win.Width(), Height: cs.win.Height(), Font: css.DefaultFontSize, Root: css.DefaultFontSize}
	if st.FontSize > 0 {
		u.Font = st.FontSize
	}
	return u
}

// u is the drawing scale, shared by every helper below.
func (cs *cssStyle) u() int { return Scale(cs.win) }

// alpha is the opacity as an alpha fraction of 255.
func (cs *cssStyle) alpha(st css.Style) int {
	if !st.Has("opacity") {
		return 255
	}
	return int(st.Opacity*255 + 0.5)
}

// ink is the text colour, faded by opacity. A widget that never declares
// colour draws with the theme's text ink; a declared colour that reverses to
// the theme sentinel does too. Only an explicit color: transparent — or any
// colour whose own alpha channel is already zero — paints nothing.
func (cs *cssStyle) ink(st css.Style) canvas.Color {
	if st.Has("color") && st.Color != css.ThemeInk {
		if st.Color.A() == 0 {
			return canvas.Transparent
		}
		return canvas.Fade(st.Color, cs.alpha(st))
	}
	return canvas.Fade(cs.win.Theme().Text, cs.alpha(st))
}

// fill is the background colour, faded by opacity. A fully transparent fill
// stays transparent, so background-color: transparent shows the theme
// instead of painting a black box.
func (cs *cssStyle) fill(st css.Style) canvas.Color {
	if st.Background.A() == 0 {
		return canvas.Transparent
	}
	return canvas.Fade(st.Background, cs.alpha(st))
}

// length measures a Length at the drawing scale against the window. A
// percentage is of the base (usually the window width), a fixed or physical
// length is its reference pixels times the scale, and the viewport/font units
// resolve against the real window and the widget's own font size.
func (cs *cssStyle) length(st css.Style, l css.Length, base int) int {
	ctx := css.Units{
		Width:  cs.win.Width(),
		Height: cs.win.Height(),
		Font:   st.FontSize,
		Root:   css.DefaultFontSize,
	}
	if l.IsPct() {
		return l.Resolve(css.Units{Width: base})
	}
	return l.Resolve(ctx) * cs.u()
}

// paintBox fills a rounded box with the background and strokes it with the
// border, honouring border-box widths, an elliptical radius and the per-side
// border styles. Around that it layers the visual effects: a backdrop-filter
// sees the pixels behind the box, the outer shadows fall behind it, the inset
// shadows fall on the surface, an outline rings it, and a filter sees the
// box once it is painted. A missing background leaves the window's own fill
// showing through.
func (cs *cssStyle) paintBox(win *antui.Window, st css.Style, x, y, w, h int) {
	cv := win.Canvas()
	rx := st.Radius[0] * cs.u()
	ry := rx
	if st.RadiusY[0] != 0 {
		ry = st.RadiusY[0] * cs.u()
	}

	if len(st.BackdropFilters) > 0 {
		cs.applyFilters(cv, st.BackdropFilters, x, y, w, h)
	}

	cs.paintShadows(win, st, x, y, w, h, rx, ry, false)

	hasBg := cs.backgroundOn(st)
	if hasBg {
		cs.paintBackground(win, st, x, y, w, h, rx, ry)
	}

	cs.paintShadows(win, st, x, y, w, h, rx, ry, true)

	if st.BorderOn() {
		cs.paintBorder(cv, st, x, y, w, h, rx, ry, cs.interior(win, st, hasBg))
	}

	if cs.outlineOn(st) {
		cs.paintOutline(win, st, x, y, w, h)
	}

	if len(st.Filters) > 0 {
		cs.applyFilters(cv, st.Filters, x, y, w, h)
	}
}

// boxOn reports whether the CSS decoration layer has anything to fill a box
// with: a background, a border, outside or inset shadows, an outline, or a
// backdrop filter. A widget whose CSS style has none of these paints its own
// face straight on the window, and the paint path skips the box.
func (cs *cssStyle) boxOn(st css.Style) bool {
	if cs.backgroundOn(st) || st.BorderOn() || cs.outlineOn(st) {
		return true
	}
	return len(st.BoxShadow) > 0 || len(st.BackdropFilters) > 0
}

// interior is the colour behind the border: the box's own fill, or the theme
// surface when the box has none.
func (cs *cssStyle) interior(win *antui.Window, st css.Style, hasBg bool) canvas.Color {
	if hasBg {
		return cs.fill(st)
	}
	return win.Theme().Surface
}

// paintBorder strokes the border. A uniform, fully-solid border is drawn as a
// thicker rounded rectangle behind the surface, which follows the radius.
// Mixed or non-solid borders are painted side by side with square corners.
func (cs *cssStyle) paintBorder(cv *canvas.Canvas, st css.Style, x, y, w, h, rx, ry int, fill canvas.Color) {
	u := cs.u()

	if !cs.perSideBorder(st) && len(st.BackgroundImages) == 0 {
		bw := st.BorderWidth[0] * u
		if bw <= 0 {
			bw = u
		}
		col := st.BoxColor[0]
		cv.FillRoundRectXY(x, y, w, h, max(rx, u), max(ry, u), canvas.Fade(col, cs.alpha(st)))
		cv.FillRoundRectXY(x+bw, y+bw, w-2*bw, h-2*bw, max(rx-u, 0), max(ry-u, 0), fill)
		return
	}
	for side := 0; side < 4; side++ {
		style := st.BoxStyle[side]
		if style == css.BorderNone {
			continue
		}
		bw := st.BorderWidth[side] * u
		if bw <= 0 {
			bw = u
		}
		cs.paintSide(cv, side, style, bw, x, y, w, h, canvas.Fade(st.BoxColor[side], cs.alpha(st)))
	}
}

// perSideBorder reports whether the four sides cannot be drawn as one ring:
// a non-solid style, a missing side, or sides of different width or colour.
func (cs *cssStyle) perSideBorder(st css.Style) bool {
	first := -1
	for i := 0; i < 4; i++ {
		if st.BoxStyle[i] == css.BorderNone {
			return true
		}
		if st.BoxStyle[i] > css.BorderSolid {
			return true
		}
		if first < 0 {
			first = i
			continue
		}
		if st.BorderWidth[i] != st.BorderWidth[first] || st.BoxColor[i] != st.BoxColor[first] {
			return true
		}
	}
	return false
}

// paintSide strokes one side of a box in the given border style. side is 0
// top, 1 right, 2 bottom, 3 left.
func (cs *cssStyle) paintSide(cv *canvas.Canvas, side int, style uint8, bw, x, y, w, h int, col canvas.Color) {
	switch style {
	case css.BorderDouble:
		cs.paintDoubleSide(cv, side, bw, x, y, w, h, col)
	case css.BorderDashed:
		cs.paintDashSide(cv, side, bw, x, y, w, h, col, false)
	case css.BorderDotted:
		cs.paintDashSide(cv, side, bw, x, y, w, h, col, true)
	default:
		cs.paintSolidSide(cv, side, bw, x, y, w, h, col)
	}
}

// paintSolidSide fills one straight side.
func (cs *cssStyle) paintSolidSide(cv *canvas.Canvas, side, bw, x, y, w, h int, col canvas.Color) {
	switch side {
	case 0:
		cv.FillRect(x, y, w, bw, col)
	case 1:
		cv.FillRect(x+w-bw, y, bw, h, col)
	case 2:
		cv.FillRect(x, y+h-bw, w, bw, col)
	default:
		cv.FillRect(x, y, bw, h, col)
	}
}

// paintDashSide draws a dashed or dotted side: dashes as rectangles, dots as
// small discs, spaced along the side.
func (cs *cssStyle) paintDashSide(cv *canvas.Canvas, side int, bw, x, y, w, h int, col canvas.Color, dotted bool) {
	horizontal := side == 0 || side == 2
	length := w
	if !horizontal {
		length = h
	}
	u := cs.u()
	step := max(3*bw, 3*u)
	if dotted {
		step = max(2*bw, 2*u)
	}

	if dotted {
		r := max(bw/2, 1)
		for p := 0; p+bw <= length; p += step {
			if horizontal {
				yy := y
				if side == 2 {
					yy = y + h - bw
				}
				cv.FillCircle(x+p+bw/2, yy+bw/2, r, col)
			} else {
				xx := x
				if side == 1 {
					xx = x + w - bw
				}
				cv.FillCircle(xx+bw/2, y+p+bw/2, r, col)
			}
		}
		return
	}

	seg := max(2*bw, 2*u)
	for p := 0; p < length; p += step {
		n := min(seg, length-p)
		if n <= 0 {
			break
		}
		if horizontal {
			yy := y
			if side == 2 {
				yy = y + h - bw
			}
			cv.FillRect(x+p, yy, n, bw, col)
		} else {
			xx := x
			if side == 1 {
				xx = x + w - bw
			}
			cv.FillRect(xx, y+p, bw, n, col)
		}
	}
}

// paintDoubleSide draws a double border: two thin lines with a gap between.
// A border too thin for a real double falls back to a solid side.
func (cs *cssStyle) paintDoubleSide(cv *canvas.Canvas, side int, bw, x, y, w, h int, col canvas.Color) {
	line := max(bw/3, 1)
	if bw-2*line < 1 {
		cs.paintSolidSide(cv, side, bw, x, y, w, h, col)
		return
	}
	switch side {
	case 0:
		cv.FillRect(x, y, w, line, col)
		cv.FillRect(x, y+bw-line, w, line, col)
	case 1:
		cv.FillRect(x+w-line, y, line, h, col)
		cv.FillRect(x+w-bw, y, line, h, col)
	case 2:
		cv.FillRect(x, y+h-line, w, line, col)
		cv.FillRect(x, y+h-bw, w, line, col)
	default:
		cv.FillRect(x, y, line, h, col)
		cv.FillRect(x+bw-line, y, line, h, col)
	}
}

// paintScrollbars draws the indicators an overflow: scroll or auto box gets
// when its content is larger than the box. They are visual for now; the
// wheel and drag behaviour arrives with the scroll model.
func (cs *cssStyle) paintScrollbars(win *antui.Window, st css.Style, x, y, w, h, contentW, contentH int) {
	u := cs.u()
	bar := max(3, 3*u)
	inset := max(1, u)
	cv := win.Canvas()
	theme := win.Theme()
	track := canvas.Fade(theme.Surface, 128)
	thumb := canvas.Shade(theme.Surface, -0.4)

	vScroll := st.Overflow[1] == css.OverflowScroll || st.Overflow[1] == css.OverflowAuto
	if vScroll && contentH > h && h > 2*inset {
		tx := x + w - bar - inset
		cv.FillRoundRect(tx, y+inset, bar, h-2*inset, bar/2, track)
		th := max(h*(h-2*inset)/max(contentH, 1), 3*bar)
		th = min(th, h-2*inset)
		cv.FillRoundRect(tx, y+inset, bar, th, bar/2, thumb)
	}
	hScroll := st.Overflow[0] == css.OverflowScroll || st.Overflow[0] == css.OverflowAuto
	if hScroll && contentW > w && w > 2*inset {
		ty := y + h - bar - inset
		cv.FillRoundRect(x+inset, ty, w-2*inset, bar, bar/2, track)
		tw := max(w*(w-2*inset)/max(contentW, 1), 3*bar)
		tw = min(tw, w-2*inset)
		cv.FillRoundRect(x+inset, ty, tw, bar, bar/2, thumb)
	}
}

// natural is the size a widget wants before any CSS width or height says
// otherwise: its text, its padding, its border, at the window's scale.
func (cs *cssStyle) natural(role, label string, st css.Style, u int) (w, h int) {
	var bx, by int
	if st.BorderOn() {
		bx = (st.BorderWidth[1] + st.BorderWidth[3]) * u
		by = (st.BorderWidth[0] + st.BorderWidth[2]) * u
	}
	labelW := cs.measure(st, label)
	switch role {
	case css.RoleLabel:
		return labelW, textHeight(u)
	case css.RoleButton:
		return labelW + 2*cs.padding(st, 3) + bx,
			textHeight(u) + cs.padding(st, 2) + cs.padding(st, 0) + by
	case css.RoleCheckbox, css.RoleRadio:
		return 18*u + 8*u + labelW, 18 * u
	case css.RoleSlider:
		return 24 * canvas.FontWidth * u, 3 * canvas.FontHeight * u
	case css.RoleInput:
		return max(labelW+2*cs.padding(st, 3)+bx, 24*canvas.FontWidth*u),
			textHeight(u) + cs.padding(st, 2) + cs.padding(st, 0) + by
	case css.RoleSelect:
		return max(labelW+2*cs.padding(st, 3)+16*u+bx, 24*canvas.FontWidth*u),
			textHeight(u) + cs.padding(st, 2) + cs.padding(st, 0) + by
	case css.RoleTextArea:
		return max(24*canvas.FontWidth*u, 32*canvas.FontWidth*u), 3 * canvas.FontHeight * u
	case css.RoleSwitch:
		return 36*u + 8*u + labelW, 20 * u
	case css.RoleProgress:
		return 24 * canvas.FontWidth * u, 2 * canvas.FontHeight * u
	case css.RoleDatePicker:
		return max(labelW+2*cs.padding(st, 3)+16*u+bx, 24*canvas.FontWidth*u),
			textHeight(u) + cs.padding(st, 2) + cs.padding(st, 0) + by
	}
	return labelW, textHeight(u)
}

// fontScale is how much the style's font-size enlarges the built-in face,
// with 16 reference pixels as the face's natural size.
func (cs *cssStyle) fontScale(st css.Style) int {
	u := cs.u()
	if !st.Has("font-size") || st.FontSize <= 0 {
		return u
	}
	scale := int(float64(st.FontSize*u)/16 + 0.5)
	if scale < 1 {
		scale = 1
	}
	return scale
}

func (cs *cssStyle) padding(st css.Style, side int) int {
	return cs.base(st, st.Padding[side])
}

// base measures a length against the window width; used for padding and
// margins.
func (cs *cssStyle) base(st css.Style, l css.Length) int {
	return cs.length(st, l, cs.win.Width())
}

// ---------------------------------------------------------------------------
// The Style interface
// ---------------------------------------------------------------------------

func (cs *cssStyle) Background(win *antui.Window) canvas.Color {
	st := cs.classes.GetStyle(css.RoleBody, nil, css.StateNone, win.Width())
	if st.Has("background-color") || st.Has("background") {
		return canvas.Fade(st.Background, cs.alpha(st))
	}
	return win.Theme().Background
}

func (cs *cssStyle) Label(win *antui.Window, x, y int, text string) {
	st := cs.styleFor(css.RoleLabel, State{})
	cs.run(win.Canvas(), st, x, y, text)
}

func (cs *cssStyle) Button(win *antui.Window, s State, x, y, w, h int, label string) {
	st := cs.styleFor(css.RoleButton, s)
	cs.paintBox(win, st, x, y, w, h)
	cs.drawText(win, st, x, y, w, h, label)
}

func (cs *cssStyle) Checkbox(win *antui.Window, s State, x, y int, label string, on bool) {
	st := cs.styleFor(css.RoleCheckbox, s)
	if on {
		st = cs.styleFor(css.RoleCheckbox, State{On: true, Focused: s.Focused, Hovered: s.Hovered})
	}
	u := cs.u()
	box := 18 * u
	cs.paintBox(win, st, x, y, box, box)
	if on {
		mark := cs.ink(st)
		cv := win.Canvas()
		cv.Line(x+4*u, y+9*u, x+7*u, y+13*u, mark)
		cv.Line(x+5*u, y+9*u, x+8*u, y+13*u, mark)
		cv.Line(x+7*u, y+13*u, x+14*u, y+5*u, mark)
		cv.Line(x+8*u, y+13*u, x+15*u, y+5*u, mark)
	}
	cs.drawText(win, st, x+box+8*u, y, box, box, label)
}

func (cs *cssStyle) Radio(win *antui.Window, s State, x, y int, label string, on bool) {
	st := cs.styleFor(css.RoleRadio, s)
	if on {
		st = cs.styleFor(css.RoleRadio, State{On: true, Focused: s.Focused, Hovered: s.Hovered})
	}
	u := cs.u()
	size := 18 * u
	half := size / 2
	cv := win.Canvas()
	cv.FillCircle(x+half, y+half, half, cs.fill(st))
	border := st.BoxColor[0]
	if !st.BorderOn() {
		border = canvas.Shade(cs.fill(st), -0.3)
	}
	cv.Circle(x+half, y+half, half, canvas.Fade(border, cs.alpha(st)))
	if on {
		cv.FillCircle(x+half, y+half, half-4*u, cs.ink(st))
	}
	cs.drawText(win, st, x+size+8*u, y, size, size, label)
}

func (cs *cssStyle) Slider(win *antui.Window, s State, x, y, w, h int, value float32) {
	st := cs.styleFor(css.RoleSlider, s)
	u := cs.u()
	value = clampFloat(value, 0, 1)
	knobRadius := max(h/2, 6*u)
	usable := max(w-2*knobRadius, 1)
	knobX := x + knobRadius + int(value*float32(usable)+0.5)
	trackY := y + h/2 - 3*u

	cv := win.Canvas()
	cv.FillRoundRect(x+knobRadius, trackY, knobX-(x+knobRadius), 6*u, 3*u, cs.fill(st))
	cv.FillRoundRect(x+knobRadius, trackY, usable, 6*u, 3*u, canvas.Shade(cs.fill(st), -0.25))
	cv.FillCircle(knobX, y+h/2, knobRadius, cs.ink(st))
	cv.FillCircle(knobX, y+h/2, knobRadius-4*u, cs.fill(st))
}

func (cs *cssStyle) Input(win *antui.Window, s State, x, y, w, h int, text string, cursor int, caret bool) {
	st := cs.styleFor(css.RoleInput, s)
	cs.paintBox(win, st, x, y, w, h)
	u := cs.u()
	padding := cs.base(st, st.Padding[3])
	savedClip := win.Canvas().Clip
	win.Canvas().SetClip(x+padding, y+1, w-2*padding, h-2)

	cursor = clampInt(cursor, 0, len(text))
	cursorPx := cs.measure(st, text[:cursor])
	scroll := max(cursorPx-(w-2*padding), 0)
	// drawText adds the style's own left padding; pass the pre-scrolled x so
	// the glyphs land at x+padding-scroll — the same origin the caret uses.
	cs.drawText(win, st, x-scroll, y, w, h, text)

	if s.Focused && caret {
		cx := x + padding - scroll + cursorPx
		cy := y + (h-cs.textHeight(st))/2
		win.Canvas().FillRect(cx, cy, u, cs.textHeight(st), cs.ink(st))
	}
	win.Canvas().Clip = savedClip
}

func (cs *cssStyle) InputTextPos(win *antui.Window, s State, x, y, w, h int, text string, cursor, mx, my int) int {
	st := cs.styleFor(css.RoleInput, s)
	padding := cs.base(st, st.Padding[3])
	cursor = clampInt(cursor, 0, len(text))
	cursorPx := cs.measure(st, text[:cursor])
	scroll := max(cursorPx-(w-2*padding), 0)
	return caretAt(func(prefix string) int { return cs.measure(st, prefix) }, text, x+padding-scroll, mx)
}

func (cs *cssStyle) TextAreaTextPos(win *antui.Window, s State, x, y, w, h int, text string, cursor, mx, my int) int {
	st := cs.styleFor(css.RoleTextArea, s)
	padding := cs.base(st, st.Padding[3])
	textW := max((w-2*padding)/cs.fontScale(st), 12)
	cursor = clampInt(cursor, 0, len(text))
	lines := wrapLines(text, textW, st.WhiteSpace, st.OverflowWrap)
	return textAreaTextPos(text, lines, padding, x, y, h, cursor, cs.lineHeight(st, cs.fontScale(st)),
		func(prefix string) int { return cs.measure(st, prefix) }, mx, my)
}

func (cs *cssStyle) textHeight(st css.Style) int {
	return canvas.TextHeight() * cs.fontScale(st)
}

func (cs *cssStyle) Select(win *antui.Window, s State, x, y, w, h int, value string, open bool) {
	st := cs.styleFor(css.RoleSelect, s)
	if open {
		st = cs.styleFor(css.RoleSelect, State{Focused: true, Hovered: s.Hovered})
	}
	cs.paintBox(win, st, x, y, w, h)
	padR := cs.base(st, st.Padding[1])
	u := cs.u()
	textW := w - padR - 14*u
	cs.drawText(win, st, x, y, max(textW, 0), h, value)

	aw, ah := 6*u, 4*u
	ax := x + w - padR - aw
	ay := y + (h-ah)/2
	cv := win.Canvas()
	if open {
		cv.Line(ax, ay+ah, ax+aw/2, ay, cs.ink(st))
		cv.Line(ax+aw/2, ay, ax+aw, ay+ah, cs.ink(st))
	} else {
		cv.Line(ax, ay, ax+aw/2, ay+ah, cs.ink(st))
		cv.Line(ax+aw/2, ay+ah, ax+aw, ay, cs.ink(st))
	}
}

func (cs *cssStyle) SelectOption(win *antui.Window, s State, x, y, w, h int, label string, selected bool) {
	st := cs.styleFor(css.RoleSelect, s)
	if selected {
		st = cs.styleFor(css.RoleSelect, State{On: true, Focused: s.Focused, Hovered: s.Hovered})
	}
	fill := cs.fill(st)
	if selected || s.Hovered {
		fill = canvas.Shade(fill, -0.15)
	}
	u := cs.u()
	cv := win.Canvas()
	cv.FillRect(x, y, w, h, canvas.Fade(fill, cs.alpha(st)))
	if selected || s.Hovered {
		cv.FillRect(x, y, 3*u, h, cs.ink(st))
	}
	cs.drawText(win, st, x, y, w, h, label)
}

func (cs *cssStyle) TextArea(win *antui.Window, s State, x, y, w, h int, text string, cursor int, caret bool) {
	st := cs.styleFor(css.RoleTextArea, s)
	cs.paintBox(win, st, x, y, w, h)
	u := cs.u()
	padding := cs.base(st, st.Padding[3])
	innerW := w - 2*padding
	textW := max(innerW/cs.fontScale(st), 12)
	cursor = clampInt(cursor, 0, len(text))
	lines := wrapLines(text, textW, st.WhiteSpace, st.OverflowWrap)

	lineH := cs.lineHeight(st, cs.fontScale(st))
	visible := max((h-2*padding)/lineH, 1)
	viewTop := 0
	for li, l := range lines {
		if cursor >= l[0] && cursor < l[1] || (li == len(lines)-1 && cursor == l[1]) {
			viewTop = max(li-visible+1, 0)
			break
		}
	}
	viewTop = max(min(viewTop, len(lines)-visible), 0)

	savedClip := win.Canvas().Clip
	win.Canvas().SetClip(x+padding, y+padding, innerW, h-2*padding)

	for lif, l := range lines {
		if lif < viewTop || lif >= viewTop+visible {
			continue
		}
		cy := y + padding + (lif-viewTop)*lineH
		hard := l[1] < len(text) && text[l[1]] == '\n'
		justify := st.TextAlign == 3 && lif < len(lines)-1 && !hard
		// drawTextLine adds the left padding itself; pass the plain x so the
		// line and the caret share the x+padding origin.
		cs.drawTextLine(win, st, x, cy, innerW, lineH, text[l[0]:l[1]], justify)
		if s.Focused && caret && cursor >= l[0] && cursor <= l[1] {
			off := cursor - l[0]
			cx := x + padding + cs.measure(st, text[l[0]:l[0]+off])
			win.Canvas().FillRect(cx, cy, u, lineH, cs.ink(st))
		}
	}
	win.Canvas().Clip = savedClip
}

func (cs *cssStyle) Switch(win *antui.Window, s State, x, y int, label string, on bool) {
	st := cs.styleFor(css.RoleSwitch, s)
	if on {
		st = cs.styleFor(css.RoleSwitch, State{On: true, Focused: s.Focused, Hovered: s.Hovered})
	}
	u := cs.u()
	trackW, trackH := 36*u, 20*u

	fill := cs.fill(st)
	if on {
		fill = cs.ink(st)
	}
	cv := win.Canvas()
	cv.FillRoundRect(x, y, trackW, trackH, trackH/2, canvas.Fade(fill, cs.alpha(st)))
	cx := x + trackH/2
	if on {
		cx = x + trackW - trackH/2
	}
	thumb := canvas.Shade(fill, -0.2)
	if on {
		thumb = cs.fill(st)
	}
	cv.FillCircle(cx, y+trackH/2, trackH/2-2*u, thumb)
	cs.run(cv, st, x+trackW+8*u, y+(trackH-cs.textHeight(st))/2, label)
}

func (cs *cssStyle) Progress(win *antui.Window, x, y, w, h int, value float32) {
	st := cs.styleFor(css.RoleProgress, State{})
	value = clampFloat(value, 0, 1)
	filled := int(float32(w)*value + 0.5)
	cv := win.Canvas()

	cv.FillRoundRect(x, y, w, h, h/2, canvas.Shade(cs.fill(st), -0.25))
	if filled > 0 {
		cv.FillRoundRect(x, y, filled, h, h/2, cs.fill(st))
	}
}

func (cs *cssStyle) DatePicker(win *antui.Window, s State, x, y, w, h int, year, month, firstWD, days, selected, today, hover int) {
	st := cs.styleFor(css.RoleDatePicker, s)
	cs.paintBox(win, st, x, y, w, h)
	u := cs.u()

	header := dateHeaderHeight(u)
	mname := monthName(month)
	labelW := textWidth(u, mname) + 3*u + textWidth(u, itoa(year))
	lx := x + (w-labelW)/2
	cv := win.Canvas()
	cs.run(cv, st, lx, y+(header-cs.textHeight(st))/2, mname)
	cs.run(cv, st, lx+textWidth(u, mname)+3*u, y+(header-cs.textHeight(st))/2, itoa(year))

	aw := 8 * u
	ay := y + (header-cs.textHeight(st))/2
	prevX, nextX := x+2*u, x+w-2*u-aw
	cv.Line(prevX+3*u, ay+cs.textHeight(st)/2, prevX+aw-3*u, ay+3*u, cs.ink(st))
	cv.Line(prevX+3*u, ay+cs.textHeight(st)/2, prevX+aw-3*u, ay+cs.textHeight(st)-3*u, cs.ink(st))
	cv.Line(nextX+3*u, ay+3*u, nextX+aw-3*u, ay+cs.textHeight(st)/2, cs.ink(st))
	cv.Line(nextX+aw-3*u, ay+cs.textHeight(st)/2, nextX+3*u, ay+cs.textHeight(st)-3*u, cs.ink(st))

	c := dateCellSize(u)
	gx, _, _, _ := dateGridRect(x, y, u)
	hy := weekdayHeadY(x, y, u)
	for d := 0; d < 7; d++ {
		cx := gx + d*c
		cs.run(cv, st, cx+(c-textWidth(u, "S"))/2, hy, dayHeadName((firstWD+d)%7))
	}

	for day := 1; day <= days; day++ {
		pos := firstWD + day - 1
		col, row := pos%7, pos/7
		cx, cy, cc, _ := dayCellRect(x, y, u, col, row)
		tx, ty := cx+(cc-textWidth(u, itoa(day)))/2, cy+(cc-cs.textHeight(st))/2
		switch day {
		case selected:
			// The selected day is a chip of the widget's own ink with the
			// number reversed out in the surface colour, the way the
			// built-in templates mark it in their accent colour.
			cv.FillRoundRect(cx, cy, cc, cc, 2*u, canvas.Fade(cs.ink(st), cs.alpha(st)))
			cs.runOn(cv, st, cs.surface(st), tx, ty, itoa(day))
		case hover:
			cv.FillRoundRect(cx, cy, cc, cc, 2*u, canvas.Fade(canvas.Shade(cs.fill(st), -0.15), cs.alpha(st)))
			cs.run(cv, st, tx, ty, itoa(day))
		default:
			cs.run(cv, st, tx, ty, itoa(day))
		}
		if day == today && day != selected {
			cv.Circle(cx+cc/2, cy+cc/2, cc/2-2*u, cs.ink(st))
		}
	}
}

// surface is the colour a reversed-out number draws in on a selected day
// chip: the widget's own background when it has one, or the theme's
// background otherwise — which, against the theme's text (the chip's ink),
// always contrasts.
func (cs *cssStyle) surface(st css.Style) canvas.Color {
	if fill := cs.fill(st); fill.A() != 0 {
		return fill
	}
	return cs.win.Theme().Background
}

func (cs *cssStyle) DatePickerBox(win *antui.Window, s State, x, y, w, h int, value string, open bool) {
	st := cs.styleFor(css.RoleDatePicker, s)
	if open {
		st = cs.styleFor(css.RoleDatePicker, State{Focused: true, Hovered: s.Hovered})
	}
	cs.paintBox(win, st, x, y, w, h)
	padR := cs.base(st, st.Padding[1])
	u := cs.u()
	textW := w - padR - 14*u
	cs.drawText(win, st, x, y, max(textW, 0), h, value)

	aw, ah := 6*u, 4*u
	ax := x + w - padR - aw
	ay := y + (h-ah)/2
	cv := win.Canvas()
	if open {
		cv.Line(ax, ay+ah, ax+aw/2, ay, cs.ink(st))
		cv.Line(ax+aw/2, ay, ax+aw, ay+ah, cs.ink(st))
	} else {
		cv.Line(ax, ay, ax+aw/2, ay+ah, cs.ink(st))
		cv.Line(ax+aw/2, ay+ah, ax+aw, ay, cs.ink(st))
	}
}

// clampFloat keeps v inside [0,1] for slider and progress values.
func clampFloat(v, lo, hi float32) float32 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
