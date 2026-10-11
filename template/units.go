package template

import "github.com/gabrielluizsf/antui/template/css"

// units is the measurement context for lengths inside a style.
func (cs *cssStyle) units(st css.Style) css.Units {
	vp := cs.viewport()
	return css.Units{
		Width:  vp.Width,
		Height: vp.Height,
		Font:   st.FontSize,
		Root:   css.DefaultFontSize,
		Scale:  vp.Scale,
		Scheme: vp.Scheme,
	}
}
