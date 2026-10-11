package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/template/audio"
)

// TemplateWithCSS makes the coord-free template: the developer passes no
// coordinates, and the template decides where each widget lives from an
// external stylesheet. Call [CSS.SetStyle] with the CSS file and a class
// table before drawing anything.
//
// The start-up loading screen comes with it: it is drawn and presented
// before this returns, so the window already wears it while the program
// loads what it needs — the stylesheet above, a font, a database. It stays
// until [CSS.Ready], and a program that never gets there is closed by the
// loading timeout. [CSS.SetDisableLoading] asks to be without it.
func TemplateWithCSS(win *antui.Window) *CSS {
	cs := newCSSStyle(win)
	c := &CSS{
		win:   win,
		style: cs,
		ui:    &uiTemplate{style: cs, sound: audio.Simple, ld: newLoading()},
	}
	c.ui.ld.showFirstFrame(win, c.ui)
	return c
}
