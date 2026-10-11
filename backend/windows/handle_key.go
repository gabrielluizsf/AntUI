//go:build windows

package windows

import "github.com/gabrielluizsf/antui/backend"

// handleKey folds a keyboard message into the window. Alt+F4 is left
// unhandled so it reaches the system, or the window cannot be closed the way
// every other window on the desktop closes.
func (d *Driver) handleKey(m wndMsg) (uintptr, bool) {
	switch m.message {
	case wmKeyDown, wmSysKeyDown:
		mods := winMods()
		d.win.Push(backend.Event{
			Type:   backend.EventKeyDown,
			Key:    winKey(m.wparam, m.lparam),
			Mods:   mods,
			Repeat: m.lparam&0x40000000 != 0,
		})
		if m.message == wmSysKeyDown && m.wparam == vkF4 && mods&backend.ModAlt != 0 {
			return 0, false
		}
		return 0, true

	case wmKeyUp, wmSysKeyUp:
		d.win.Push(backend.Event{
			Type: backend.EventKeyUp,
			Key:  winKey(m.wparam, m.lparam),
			Mods: winMods(),
		})
		return 0, true

	case wmChar, wmSysChar:
		return d.handleChar(uint16(m.wparam)), true
	}
	return 0, false
}
