//go:build android

package antui

import (
	"github.com/gabrielluizsf/antui/backend/android/ndk"
	"github.com/gabrielluizsf/antui/canvas"
)

// present blits the frame into the buffer the compositor lends us, reporting
// whether it got posted.
//
// The rectangle handed in is what changed. The compositor's answer to it is
// usually larger — it hands out buffers in rotation, so the one being locked
// holds a frame from two or three ago, and everything stale in it has to be
// written. Lock rewrites the rectangle with what must be filled, and that is
// the one to blit.
func (b *androidWindow) present(win *Window, dirty canvas.Area) bool {
	if b.surface == nil || dirty.Width <= 0 || dirty.Height <= 0 {
		return false
	}
	r := ndk.Rect{
		Left:   int32(dirty.X),
		Top:    int32(dirty.Y),
		Right:  int32(dirty.X + dirty.Width),
		Bottom: int32(dirty.Y + dirty.Height),
	}
	buf, err := b.surface.Lock(&r)
	if err != nil {
		// The surface can be taken away between the event and this call.
		// Losing a frame is the right answer; the WindowDown event is on its
		// way and pump will park.
		ndk.Warnf("cannot lock the surface: %v", err)
		b.surface = nil
		return false
	}
	blit(buf, win.cv, r)
	if err := b.surface.UnlockAndPost(); err != nil {
		ndk.Warnf("cannot post the frame: %v", err)
		b.surface = nil
		return false
	}
	return true
}
