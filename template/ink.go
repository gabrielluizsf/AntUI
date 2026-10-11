package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

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
