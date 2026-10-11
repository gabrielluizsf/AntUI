//go:build windows

package windows

import "github.com/gabrielluizsf/antui/backend"

// mouseButton answers a button press or release. Capturing on a left press
// means a drag that leaves the window still reports its release, so a button
// or slider cannot get stuck held down.
func (d *Driver) mouseButton(m wndMsg) uintptr {
	var t backend.EventType
	var button backend.MouseButton
	switch m.message {
	case wmLButtonDown:
		procSetCapture.Call(m.hwnd)
		t, button = backend.EventMouseDown, backend.MouseLeft
	case wmLButtonUp:
		procReleaseCapture.Call()
		t, button = backend.EventMouseUp, backend.MouseLeft
	case wmRButtonDown:
		t, button = backend.EventMouseDown, backend.MouseRight
	case wmRButtonUp:
		t, button = backend.EventMouseUp, backend.MouseRight
	case wmMButtonDown:
		t, button = backend.EventMouseDown, backend.MouseMiddle
	case wmMButtonUp:
		t, button = backend.EventMouseUp, backend.MouseMiddle
	}
	d.win.Push(backend.Event{
		Type:   t,
		Button: button,
		X:      loWord(m.lparam),
		Y:      hiWord(m.lparam),
		Mods:   winMods(),
	})
	return 0
}
