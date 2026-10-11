package template

import "github.com/gabrielluizsf/antui"

// wipe puts the plain background where the screen was, and presents it. It
// is for the developer who turns the screen off after it already appeared —
// the flag set once the template exists — so that a program which asked to
// be without the loading picture does not sit on one until its first frame.
func (l *loading) wipe(win *antui.Window, painter loadingPainter) {
	if win == nil || l == nil || !l.pushed.Load() || l.ready.Load() {
		return
	}
	win.Clear(painter.Background(win))
	win.End()
}
