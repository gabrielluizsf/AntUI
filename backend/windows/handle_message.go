//go:build windows

package windows

import "github.com/gabrielluizsf/antui/backend"

// handleMessage folds one window message into the window, reporting whether
// the message was one this package handles and what the window procedure
// should return for it. An unhandled message falls through to Windows.
func (d *Driver) handleMessage(m wndMsg) (uintptr, bool) {
	win := d.win
	switch m.message {
	case wmDropFiles:
		d.dropped(m.wparam)
		return 0, true

	case wmClose:
		win.PushSimple(backend.EventClose)
		return 0, true

	case wmDestroy:
		win.SetShouldClose()
		return 0, true

	case wmEraseBkgnd:
		return 1, true // the frame covers everything, so erasing only flickers

	case wmSettingChange, wmThemeChanged:
		// The user flipped a theme somewhere in Settings. The registry is
		// not read here — the frame loop asks when it next draws — so this
		// only says the answer it has is out of date.
		d.touchTheme()
		return 0, true

	case wmPaint:
		d.paint(m.hwnd)
		return 0, true

	case wmSize:
		d.resize(win, m.wparam, m.lparam)
		return 0, true

	case wmMouseMove, wmMouseLeave, wmLButtonDown, wmLButtonUp,
		wmRButtonDown, wmRButtonUp, wmMButtonDown, wmMButtonUp, wmMouseWheel:
		return d.handleMouse(m)

	case wmKeyDown, wmSysKeyDown, wmKeyUp, wmSysKeyUp, wmChar, wmSysChar:
		return d.handleKey(m)

	case wmSetFocus, wmKillFocus:
		win.Push(backend.Event{Type: backend.EventFocus, Focused: m.message == wmSetFocus})
		return 0, true

	case wmGetMinMaxInfo, wmSizing:
		return d.handleResizeMessage(m)

	case wmSetCursor:
		return d.handleSetCursor(m)
	}
	return 0, false
}
