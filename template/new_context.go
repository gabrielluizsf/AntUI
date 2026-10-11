package template

import "github.com/gabrielluizsf/antui"

// NewContext starts a session on a window with a template and nothing in
// view. Run, or Draw once a frame, shows the first screen.
//
// The loading screen comes with it: the template's start-up screen is drawn
// and presented right here, before the program has loaded anything, and
// stays until Ready says the program is up — or the loading timeout closes
// the app. A program that called SetDisableLoading on the template before
// handing it over never shows it.
func NewContext(win *antui.Window, tpl Template) *Context {
	c := &Context{win: win, tpl: tpl}
	if tpl != nil {
		tpl.loadingState().showFirstFrame(win, tpl)
	}
	return c
}
