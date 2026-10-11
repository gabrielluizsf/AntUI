package template

import "github.com/gabrielluizsf/antui"

// arm puts the screen on the frame being painted, by registering it to draw
// at the very end — [antui.Window.WidgetFrontmost] — above whatever the
// screen underneath asked for. Arm is called at the START of the frame's
// drawing, and Window.End runs the registered draws in reverse, so the
// loading screen lands last of all: above a dropdown, above a popup, above
// anything the app drew while it was still loading.
func (l *loading) arm(win *antui.Window, painter loadingPainter) {
	if win == nil || !l.active() {
		return
	}
	win.WidgetFrontmost(func() { l.draw(win, painter) })
}
