//go:build android

package antui

import (
	"errors"
	"github.com/gabrielluizsf/antui/backend/android/app"
)

// open waits for the system to give the app a surface, and takes its size
// from it. The width and height asked for are ignored, because there is
// nothing to ask: a phone's window is the screen.
func (b *androidWindow) open(win *Window, title string, width, height int) error {
	b.activity = app.Current()
	if b.activity == nil {
		return errNotAndroidApp
	}
	b.scale = scaleFromDensity(b.activity.Density())
	b.dark = b.activity.Config().Night
	// There is no other kind of input here, so a game does not have to wait
	// for a finger to know it needs controls on the screen.
	win.SetTouchFirst()

	// The activity exists before its surface does. Everything up to the
	// first WindowUp is handled here rather than being dropped, so that a
	// Start and a Resume that arrive first are not lost.
	for b.surface == nil {
		e, ok := b.activity.Next()
		if !ok {
			return errors.New("antui: the activity was destroyed before it had a surface")
		}
		b.handle(win, e)
		if win.shouldClose {
			return errors.New("antui: the activity closed before it had a surface")
		}
	}
	b.activity.Release()
	b.sync(win)
	return nil
}
