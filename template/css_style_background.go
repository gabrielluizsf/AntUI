package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

func (cs *cssStyle) Background(win *antui.Window) canvas.Color {
	st := cs.classes.GetStyleViewport(css.StyleQuery{
		Tag:      css.RoleBody,
		State:    css.StateNone,
		Viewport: cs.viewport(),
	})
	if st.Has("background-color") || st.Has("background") {
		return fadeColor(st.Background, cs.alpha(st))
	}
	return win.Theme().Background
}
