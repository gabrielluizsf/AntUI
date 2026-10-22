//go:build windows

package windows

import "github.com/gabrielluizsf/antui/backend"

// Driver is one window on Windows. It owns the platform connections the way
// the window core owns the input state.
type Driver struct {
	win  backend.Face
	hwnd uintptr
	// The icons the window is showing, kept so they can be let go of when a
	// new one takes their place.
	iconBig, iconSmall uintptr
	instance           uintptr

	// cursor is the custom cursor [Driver.SetCursor] last gave the window,
	// reapplied every time Windows asks with WM_SETCURSOR — it forgets the
	// pointer's picture on its own the moment the mouse so much as moves.
	// Zero means the ordinary system arrow.
	cursor uintptr

	limits        backend.Limits
	info          bitmapInfoHeader
	highSurrogate uint16 // the high half of a pending UTF-16 pair
	trackingMouse bool
	classAtom     uintptr

	fullscreen                    bool
	windowedWidth, windowedHeight int // the size to come back out to
	mouseX, mouseY                int // the last position the pointer reported

	// The color scheme Windows paints its own apps in, read from the
	// registry by SystemDark in theme.go and read again once Windows posts
	// the setting change that says it moved.
	themeDark, themeKnown, themeRead, themeStale bool
}
