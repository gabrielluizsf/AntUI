//go:build windows

package windows

import (
	"github.com/gabrielluizsf/antui/backend"
)

// SetLimits keeps the limits for the two messages that answer with them, and
// rewrites the style bits, which is what the frame and the maximise button
// are made of.
func (d *Driver) SetLimits(limits backend.Limits) {
	d.limits = limits
	if d.hwnd == 0 || d.fullscreen {
		// A full screen is a popup with no frame at all, and putting the
		// windowed style bits back now would draw one over the game.
		return
	}

	style, _, _ := procGetWindowLongPtr.Call(d.hwnd, gwlStyle)
	wanted := style
	if limits.Fixed {
		wanted &^= uintptr(wsThickFrame | wsMaximizeBox)
	} else {
		wanted |= wsThickFrame
		if limits.NoMaximize {
			wanted &^= uintptr(wsMaximizeBox)
		} else {
			wanted |= wsMaximizeBox
		}
	}
	if wanted == style {
		return
	}
	procSetWindowLongPtr.Call(d.hwnd, gwlStyle, wanted)
	// Windows caches the frame it drew; without SWP_FRAMECHANGED the button
	// stays on screen until something else happens to repaint it.
	procSetWindowPos.Call(d.hwnd, hwndTop, 0, 0, 0, 0,
		swpNoMove|swpNoZOrder|swpNoActivate|swpFrameChanged)
}
