package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// outlineOn reports whether the outline paints. An outline with no width or
// style none is not drawn, the way a border with no width is not.
func (cs *cssStyle) outlineOn(st css.Style) bool {
	return st.OutlineWidth > 0 && st.OutlineStyle != css.BorderNone
}

// paintOutline rings the box outside its border, pushed out by outline-offset.
// The ring is drawn side by side, so its corners are square rather than mitered.
func (cs *cssStyle) paintOutline(win *antui.Window, st css.Style, x, y, w, h int) {
	u := cs.u()
	off := st.OutlineOffset * u
	ow := st.OutlineWidth * u
	col := cs.effectColor(win, st, st.OutlineColor)
	ox, oy := x-off-ow, y-off-ow
	ow2, oh2 := w+2*(off+ow), h+2*(off+ow)
	for side := 0; side < 4; side++ {
		cs.paintSide(win.Canvas(), side, st.OutlineStyle, ow, ox, oy, ow2, oh2, col)
	}
}

// paintShadows paints the box's shadows on the given side of the surface. A
// drop-shadow filter is painted with the outer shadows, since both fall
// behind the box.
func (cs *cssStyle) paintShadows(win *antui.Window, st css.Style, x, y, w, h, rx, ry int, inset bool) {
	if !inset {
		for i := range st.Filters {
			if st.Filters[i].Drop != nil {
				cs.paintShadow(win, st, *st.Filters[i].Drop, x, y, w, h, rx, ry, false)
			}
		}
	}
	for i := range st.BoxShadow {
		if st.BoxShadow[i].Inset == inset {
			cs.paintShadow(win, st, st.BoxShadow[i], x, y, w, h, rx, ry, inset)
		}
	}
}

// paintShadow draws one shadow, scaling its lengths with the window.
func (cs *cssStyle) paintShadow(win *antui.Window, st css.Style, sh css.Shadow, x, y, w, h, rx, ry int, inset bool) {
	u := cs.u()
	col := cs.effectColor(win, st, sh.Color)
	win.Canvas().BoxShadow(x, y, w, h, rx, ry, sh.X*u, sh.Y*u, sh.Blur*u, sh.Spread*u, col, inset)
}

// effectColor turns a computed effect colour into pixels: a zero colour is the
// theme ink, and every colour is faded by the element's opacity.
func (cs *cssStyle) effectColor(win *antui.Window, st css.Style, c canvas.Color) canvas.Color {
	if c == 0 {
		c = win.Theme().Text
	}
	return canvas.Fade(c, cs.alpha(st))
}

// applyFilters runs a filter list over a rectangle of the window. A blur
// length is scaled with the window; drop-shadow is left to paintShadows.
func (cs *cssStyle) applyFilters(cv *canvas.Canvas, filters []css.Filter, x, y, w, h int) {
	u := cs.u()
	for _, f := range filters {
		if f.Drop != nil {
			continue
		}
		amount := f.Amount
		if f.Kind == canvas.FilterBlur {
			amount = f.Amount * float64(u)
		}
		cv.FilterRegion(x, y, w, h, f.Kind, amount)
	}
}

// paintTextShadows draws the text-shadow list under a text run. A shadow with
// no blur is the glyphs painted again, offset; a blurred one is rendered to a
// small layer so the glyph edges can soften.
func (cs *cssStyle) paintTextShadows(win *antui.Window, st css.Style, tx, ty int, text string, o canvas.TextStyle) {
	if text == "" {
		return
	}
	cv := win.Canvas()
	u := cs.u()
	for i := range st.TextShadow {
		sh := st.TextShadow[i]
		col := cs.effectColor(win, st, sh.Color)
		dx, dy := sh.X*u, sh.Y*u
		blur := sh.Blur * u
		if blur <= 0 {
			cv.DrawStyled(tx+dx, ty+dy, text, col, o)
			continue
		}
		pad := blur + 2
		layer, err := canvas.NewCanvas(cv.TextWidthStyled(text, o)+2*pad, canvas.TextHeight()*o.Scale+2*pad)
		if err != nil {
			continue
		}
		layer.DrawStyled(pad, pad, text, col|0xFF000000, o)
		layer.Blur(0, 0, layer.Width, layer.Height, max(blur/2, 1))
		canvas.ScaleAlpha(layer, float64(col.A())/255)
		cv.BlitOver(tx+dx-pad, ty+dy-pad, layer)
	}
}
