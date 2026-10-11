package template

import "github.com/gabrielluizsf/antui"

// showFirstFrame paints the screen the moment the template is made and puts
// it on the display, before the program has loaded anything. The window
// shows nothing until the first frame of the loop is presented, and
// everything a program does before that loop — reading a stylesheet, opening
// a database, fetching a picture — is exactly the stretch of time the screen
// is for, so the screen has to come before the loop rather than in it.
func (l *loading) showFirstFrame(win *antui.Window, painter loadingPainter) {
	if win == nil || !l.active() || l.pushed.Load() {
		return
	}
	l.pushed.Store(true)
	l.draw(win, painter)
	win.End()
}
