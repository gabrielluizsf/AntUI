package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

func (cs *cssStyle) Background(win *antui.Window) canvas.Color {
	st := cs.classes.GetStyleViewport(css.RoleBody, nil, css.StateNone, cs.viewport())
	if st.Has("background-color") || st.Has("background") {
		return canvas.Fade(st.Background, cs.alpha(st))
	}
	return win.Theme().Background
}
