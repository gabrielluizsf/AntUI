package template

import "github.com/gabrielluizsf/antui/template/css"

// unitsCtx is the measurement context an animated style resolves lengths
// and viewport units in: the window and the widget's own font size.
func (cs *cssStyle) unitsCtx(st css.Style) css.Units {
	vp := cs.viewport()
	u := css.Units{
		Width: vp.Width, Height: vp.Height,
		Font: css.DefaultFontSize, Root: css.DefaultFontSize,
		Scale: vp.Scale, Scheme: vp.Scheme,
	}
	if st.FontSize > 0 {
		u.Font = st.FontSize
	}
	return u
}
