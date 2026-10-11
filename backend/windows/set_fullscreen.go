//go:build windows

package windows

// SetFullscreen takes the border off and covers the monitor.
//
// Windows has no full-screen state to ask for: full screen is a window with
// no border covering the monitor. So the border comes off, the window moves
// over the monitor it is mostly on — not the primary one, which is the bug
// every second implementation of this has — and coming back out puts the
// style and the placement back.
func (d *Driver) SetFullscreen(on bool) bool {
	if d.hwnd == 0 {
		return false
	}
	style, _, _ := procGetWindowLongPtr.Call(d.hwnd, gwlStyle)
	if on {
		return d.enterFullscreen(style)
	}
	return d.leaveFullscreen(style)
}
