//go:build android

package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// safe passes the system's content rectangle on as the window's safe area.
//
// It comes from the platform for nothing — onContentRectChanged is a
// callback, not a call — which is why it is used rather than the window
// insets, which are a Java round trip. The insets say more (they separate
// the keyboard and the cutout) and app.WindowInsets is there for a program
// that needs the difference.
func (b *androidWindow) safe(win *Window) {
	r := b.activity.ContentRect()
	if r.Right <= r.Left || r.Bottom <= r.Top {
		return
	}
	win.SetSafeArea(canvas.Area{
		X:      int(r.Left),
		Y:      int(r.Top),
		Width:  int(r.Right - r.Left),
		Height: int(r.Bottom - r.Top),
	})
}
