package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// fill is the background colour, faded by opacity. A fully transparent fill
// stays transparent, so background-color: transparent shows the theme
// instead of painting a black box.
func (cs *cssStyle) fill(st css.Style) canvas.Color {
	if st.Background.A() == 0 {
		return canvas.Transparent
	}
	return fadeColor(st.Background, cs.alpha(st))
}
