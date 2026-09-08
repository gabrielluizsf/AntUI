//go:build darwin && !cgo

package macos

import (
	"errors"

	"github.com/gabrielluizsf/antui/backend"
	"github.com/gabrielluizsf/antui/canvas"
)

// AppKit is Objective-C, so the macOS backend needs cgo. Without it there is
// no window to open, and saying so plainly beats an undefined symbol. The
// Canvas API works everywhere, window or not, so the library builds and runs
// with CGO_ENABLED=0 as long as no window is opened.

// Driver is one window on macos.
type Driver struct{}

var errNoCgo = errors.New("antui: opening a window on macOS needs cgo, which this build has turned off (set CGO_ENABLED=1); the Canvas API works without it")

// Open creates the window that Options describes.
func Open(win backend.Face, opts backend.Options) (*Driver, error) {
	if win == nil {
		return nil, errors.New("antui: macos backend with no window to drive")
	}
	return &Driver{}, errNoCgo
}

// Close destroys the window and lets go of everything the platform held.
func (d *Driver) Close() {}

// Pump handles whatever happened this frame, folding it into the window
// through win.push.
func (d *Driver) Pump(win backend.Face) {}

// Present copies the window's dirty rectangle to the platform surface.
func (d *Driver) Present(win backend.Face, dirty canvas.Area) {}

// SetTitle names the window in the system's decoration.
func (d *Driver) SetTitle(title string) {}

// SetFullscreen asks the system for a full-screen window, reporting whether
// the request could be made at all.
func (d *Driver) SetFullscreen(on bool) bool { return false }

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

// Clipboard is what the system clipboard holds.
func (d *Driver) Clipboard() (text string, files []string) { return "", nil }

// SetIcon gives the window an icon where the system shows one.
func (d *Driver) SetIcon(images []*canvas.Canvas) bool { return false }

// SetClipboard puts text on the system clipboard, reporting whether the
// system took it.
func (d *Driver) SetClipboard(text string) bool { return false }