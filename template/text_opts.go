package template

import (
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
	if ff := cs.fontFace(st); ff != nil {
		// The file already leans or already carries the weight it was cut
		// for, so only the part it does not have is faked on top of it.
		o.Face = ff.Face()
		o.Bold = o.Bold && !ff.Bold()
		o.Italic = o.Italic && !ff.Slanted()
	}
	return o
}
