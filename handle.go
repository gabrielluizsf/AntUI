//go:build android

package antui

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

// handle turns one activity event into what the window core understands.
func (b *androidWindow) handle(win *Window, e app.Event) {
	switch e {
	case app.WindowUp:
		b.surface = b.activity.Window()
		if b.surface != nil {
			// Ask for eight bits a channel rather than take what the surface
			// offers: a device that defaults to 565 would otherwise quietly
			// halve the colour depth of everything drawn.
			if err := b.surface.SetGeometry(0, 0, ndk.RGBA8888); err != nil {
				ndk.Warnf("cannot set the surface format: %v", err)
			}
			b.format = b.surface.Format()
		}
		win.PushSimple(EventExpose)

	case app.WindowDown:
		// The window is already gone from the activity; dropping it here as
		// well is what makes present a no-op instead of a crash.
		b.surface = nil

	case app.Redraw, app.WindowResized:
		win.PushSimple(EventExpose)

	case app.FocusGained:
		win.Push(Event{Type: EventFocus, Focused: true})
	case app.FocusLost:
		win.Push(Event{Type: EventFocus, Focused: false})

	case app.ContentRect:
		b.safe(win)

	case app.ConfigChanged:
		b.scale = scaleFromDensity(b.activity.Density())
		b.dark = b.activity.Config().Night
		b.safe(win)
		win.PushSimple(EventExpose)

	case app.Destroy:
		b.surface = nil
		win.Push(Event{Type: EventClose})
	}
}
