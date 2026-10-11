package template

import "github.com/gabrielluizsf/antui/template/css"

// length measures a Length at the drawing scale against the window. A
// percentage is of the base (usually the window width), a fixed or physical
// length is its reference pixels times the scale, and the viewport/font units
// resolve against this frame's window and the widget's own font size.
func (cs *cssStyle) length(st css.Style, l css.Length, base int) int {
	vp := cs.viewport()
	ctx := css.Units{
		Width:  vp.Width,
		Height: vp.Height,
		Font:   st.FontSize,
		Root:   css.DefaultFontSize,
		Scale:  vp.Scale,
		Scheme: vp.Scheme,
	}
	if l.IsPct() {
		return l.Resolve(css.Units{Width: base})
	}
	return l.Resolve(ctx) * cs.u()
}
