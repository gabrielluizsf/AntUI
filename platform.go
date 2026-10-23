package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// backend is what each platform implements. The window core above it knows
// nothing about X11, Win32 or Cocoa; it hands the backend a finished frame
// and asks it for events.
//
// It lives in this package and not in a backend package of its own because it
// is both handed a *Window to pump events into and called "on" that Window —
// the two halves are mutually recursive by design. A platform's work lives in
// backend/linux, backend/macos and backend/windows; the
// thin backend_*.go files in this package sit between them and this
// interface, turning each platform package's exported calls into what Window
// above expects.
type platform interface {
	open(win *Window, title string, width, height int) error
	close()
	pump(win *Window)
	// present reports whether the dirty rectangle actually reached the
	// display, so a connection that drops partway through a frame does not
	// get mistaken for one that delivered it.
	present(win *Window, dirty canvas.Area) bool
	setTitle(title string)
	// setFullscreen reports whether the request could be made at all, not
	// whether the window manager honoured it.
	setFullscreen(on bool) bool
	// setOpacity fades the whole window — decorations and all — toward
	// transparent, so what sits behind the program shows through. 255 is
	// fully opaque. It reports whether the request could be made, not whether
	// a compositor honoured it.
	setOpacity(alpha uint8) bool
	displaySize() (w, h int, ok bool)
	displayRefresh() int
	// setLimits passes the size constraints on to the window manager. They
	// are advice on X11 and enforced on Win32 and Cocoa; see window_size.go.
	setLimits(win *Window, limits Limits)
	// setSize asks for a new drawable size, reporting whether the request
	// could be made rather than whether it was honoured.
	setSize(win *Window, width, height int) bool
	// displayScale is the display's scale, and whether it said one: the raw
	// answer behind Window.DisplayScale, read once at open and again every
	// frame so a scale the display changed reaches the frame after it. See
	// window_system.go.
	displayScale() (float64, bool)
	// systemDark is the color scheme the system paints its own interface in,
	// and whether the system said at all — a platform with no way to ask
	// answers known=false rather than guessing. See window_system.go.
	systemDark() (dark, known bool)
	// clipboard is what the system clipboard holds: its text, and the files
	// it names when it names any. Both empty when there is nothing there or
	// the platform cannot be asked.
	clipboard() (text string, files []string)
	// setIcon gives the window an icon, largest first. It reports whether the
	// window system took it.
	setIcon(images []*canvas.Canvas) bool
	// setCursor replaces the mouse pointer with icon's picture, reporting
	// whether the window system took it.
	setCursor(icon CursorIcon) bool
	// sendNotification shows a system notification — the desktop's toast,
	// the phone's shade — reporting whether the system took it, not whether
	// the user saw it.
	sendNotification(title, body string) bool
	// setClipboard puts text on the system clipboard, reporting whether the
	// system took it.
	setClipboard(text string) bool
}
