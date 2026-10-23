//go:build darwin && !cgo

package macos

import (
	"github.com/gabrielluizsf/antui/backend"
	"github.com/gabrielluizsf/antui/canvas"
)

// SetOpacity fades the whole window so what sits behind it shows through.
// There is no window to fade here.
func (d *Driver) SetOpacity(alpha uint8) bool { return false }

// DisplaySize is the size of the display the window is on.
func (d *Driver) DisplaySize() (w, h int, ok bool) { return 0, 0, false }

// DisplayRefresh is how often the display repaints, in hertz.
func (d *Driver) DisplayRefresh() int { return 0 }

// SetLimits passes the size constraints on to the platform.
func (d *Driver) SetLimits(limits backend.Limits) {}

// SetSize asks the platform for a new drawable size, reporting whether the
// request could be made.
func (d *Driver) SetSize(win backend.Face, width, height int) bool { return false }

// ContentScale is pixels per point, or 0 when the system does not say.
func (d *Driver) ContentScale() float64 { return 0 }

// SystemDark is the appearance the system paints in, and without cgo there is
// no AppKit to ask — the answer is that there is no answer.
func (d *Driver) SystemDark() (bool, bool) { return false, false }

// Clipboard is what the system clipboard holds.
func (d *Driver) Clipboard() (text string, files []string) { return "", nil }

// SetIcon gives the window an icon where the system shows one.
func (d *Driver) SetIcon(images []*canvas.Canvas) bool { return false }

// SetCursor replaces the mouse pointer with icon's picture. There is no
// window to put it on here.
func (d *Driver) SetCursor(icon backend.CursorIcon) bool { return false }

// SendNotification shows a notification, which without cgo has no
// UserNotifications to go through — osascript is the one way left.
func (d *Driver) SendNotification(title, body string) bool {
	return osascriptNotify(title, body) == nil
}

// SetClipboard puts text on the system clipboard, reporting whether the
// system took it.
func (d *Driver) SetClipboard(text string) bool { return false }
