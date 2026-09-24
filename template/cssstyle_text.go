package template

// The CSS text model: turning the computed typography properties into canvas
// draw calls. Weight, slant and spacing ride in a canvas.TextStyle; transform,
// alignment, ellipsis, line-height, vertical-align and decoration are worked
// out here around the run itself.

import (
	"strings"
	"unicode"

	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// textOpts translates a computed style into the canvas text style at the
// window's scale.
func (cs *cssStyle) textOpts(st css.Style) canvas.TextStyle {
	o := canvas.TextStyle{
		Scale:         cs.fontScale(st),
		Bold:          st.FontWeight >= 600,
		Italic:        st.FontStyle != css.FontStyleNormal,
		LetterSpacing: cs.length(st, st.LetterSpacing, cs.win.Width()),
		WordSpacing:   cs.length(st, st.WordSpacing, cs.win.Width()),
	}
	return o
}

// measure is how wide text is once the style's transform, weight, slant and
// spacing have been applied.
func (cs *cssStyle) measure(st css.Style, text string) int {
	o := cs.textOpts(st)
	return cs.win.Canvas().TextWidthStyled(transformText(text, st.TextTransform), o)
}

// run draws one styled run with no alignment, for widgets that place their own
// text (switches, calendar cells).
func (cs *cssStyle) run(cv *canvas.Canvas, st css.Style, x, y int, text string) {
	o := cs.textOpts(st)
	cv.DrawStyled(x, y, transformText(text, st.TextTransform), cs.ink(st), o)
}

// runOn draws one styled run in a given colour, for a cell that paints its
// text against a fill of its own choosing (a calendar's selected day).
func (cs *cssStyle) runOn(cv *canvas.Canvas, st css.Style, colour canvas.Color, x, y int, text string) {
	o := cs.textOpts(st)
	cv.DrawStyled(x, y, transformText(text, st.TextTransform), colour, o)
}

// lineHeight is the distance between two baselines: the font's own height
// enlarged by line-height, or multiplied by it when the style set one.
func (cs *cssStyle) lineHeight(st css.Style, scale int) int {
	base := canvas.TextHeight() * scale
	if st.LineHeight > 0 {
		return max(int(float64(base)*st.LineHeight+0.5), 1)
	}
	return base
}

// vAlignShift moves the text top by vertical-align within its box.
func (cs *cssStyle) vAlignShift(st css.Style, th, h, padT, padB int) int {
	center := (h - padT - padB - th) / 2
	switch st.VerticalAlign {
	case css.VerticalAlignSub:
		return th / 4
	case css.VerticalAlignSuper:
		return -th / 4
	case css.VerticalAlignMiddle:
		return -th / 8
	case css.VerticalAlignTop, css.VerticalAlignTextTop:
		return -center
	case css.VerticalAlignBottom, css.VerticalAlignTextBottom:
		return center
	}
	if !st.BaselineShift.Auto() && !st.BaselineShift.None() {
		return -st.BaselineShift.Resolve(cs.units(st)) * cs.u()
	}
	return 0
}

// units is the measurement context for lengths inside a style.
func (cs *cssStyle) units(st css.Style) css.Units {
	return css.Units{
		Width:  cs.win.Width(),
		Height: cs.win.Height(),
		Font:   st.FontSize,
		Root:   css.DefaultFontSize,
	}
}

// fitText truncates a run with an ellipsis so it fits maxW, working rune by
// rune so the cut never splits a character.
func fitText(cv *canvas.Canvas, text string, o canvas.TextStyle, maxW int) string {
	if maxW <= 0 || cv.TextWidthStyled(text, o) <= maxW {
		return text
	}
	const dot = "…"
	dw := cv.TextWidthStyled(dot, o)
	if dw > maxW {
		return ""
	}
	runes := []rune(text)
	for i := len(runes) - 1; i > 0; i-- {
		if cv.TextWidthStyled(string(runes[:i])+dot, o) <= maxW {
			return string(runes[:i]) + dot
		}
	}
	return dot
}

// justifyGap is the extra space each gap takes when a line is justified to
// fill avail, or zero when it cannot grow.
func justifyGap(cv *canvas.Canvas, text string, o canvas.TextStyle, avail int) int {
	gaps := strings.Count(text, " ")
	if gaps == 0 {
		return 0
	}
	if extra := (avail - cv.TextWidthStyled(text, o)) / gaps; extra > 0 {
		return extra
	}
	return 0
}

// decorate paints underline, line-through and overline across a run.
func (cs *cssStyle) decorate(cv *canvas.Canvas, st css.Style, x, y, tw, th, scale int) {
	if st.TextDecoration == 0 || tw <= 0 {
		return
	}
	ink := cs.ink(st)
	u := max(scale, 1)
	if st.TextDecoration&css.TextDecorationUnderline != 0 {
		cv.FillRect(x, y+th-2*u, tw, u, ink)
	}
	if st.TextDecoration&css.TextDecorationLineThrough != 0 {
		cv.FillRect(x, y+th/2, tw, u, ink)
	}
	if st.TextDecoration&css.TextDecorationOverline != 0 {
		cv.FillRect(x, y, tw, u, ink)
	}
}

// transformText applies text-transform to a run.
func transformText(s string, t uint8) string {
	switch t {
	case css.TextTransformUppercase:
		return strings.ToUpper(s)
	case css.TextTransformLowercase:
		return strings.ToLower(s)
	case css.TextTransformCapitalize:
		return capitalize(s)
	}
	return s
}

// capitalize uppercases the first letter of every word.
func capitalize(s string) string {
	var b strings.Builder
	start := true
	for _, r := range s {
		if unicode.IsLetter(r) && start {
			b.WriteRune(unicode.ToUpper(r))
			start = false
			continue
		}
		start = !unicode.IsLetter(r) && !unicode.IsDigit(r)
		b.WriteRune(r)
	}
	return b.String()
}

// drawText paints a text run inside a box, aligned by text-align. Justifying
// belongs to wrapped paragraphs, so a lone run is never stretched.
func (cs *cssStyle) drawText(win *antui.Window, st css.Style, x, y, w, h int, text string) {
	cs.drawTextLine(win, st, x, y, w, h, text, false)
}

func (cs *cssStyle) drawTextLine(win *antui.Window, st css.Style, x, y, w, h int, text string, justify bool) {
	cv := win.Canvas()
	o := cs.textOpts(st)
	text = transformText(text, st.TextTransform)

	padL, padR := cs.padding(st, 3), cs.padding(st, 1)
	padT, padB := cs.padding(st, 2), cs.padding(st, 0)
	avail := w - padL - padR
	th := canvas.TextHeight() * o.Scale

	if justify {
		o.WordSpacing += justifyGap(cv, text, o, avail)
	}
	if st.TextOverflow == css.TextOverflowEllipsis {
		text = fitText(cv, text, o, avail)
	}
	tw := cv.TextWidthStyled(text, o)
	tx := x + padL
	if !justify {
		switch st.TextAlign {
		case 1:
			tx = x + w/2 - tw/2
		case 2:
			tx = x + w - padL - tw
		}
	}
	ty := y + padT + (h-padT-padB-th)/2 + cs.vAlignShift(st, th, h, padT, padB)
	if len(st.TextShadow) > 0 {
		cs.paintTextShadows(win, st, tx, ty, text, o)
	}
	cv.DrawStyled(tx, ty, text, cs.ink(st), o)
	cs.decorate(cv, st, tx, ty, tw, th, o.Scale)
}
