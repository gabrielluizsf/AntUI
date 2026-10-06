//go:build (linux && !android) || freebsd || openbsd || netbsd || dragonfly

package antui

import (
	"time"

	"github.com/gabrielluizsf/antui/backend"
	x11 "github.com/gabrielluizsf/antui/backend/linux"
	"github.com/gabrielluizsf/antui/canvas"
)

// x11Native is the Linux backend, in two halves: this struct adapts the
// window core's calls to the platform package, and the package itself lives
// in backend/linux.
type x11Native struct{ d *x11.Driver }

// newBackend is what OpenWith asks for when there is not yet a window.
func newBackend() platform { return &x11Native{} }

func (n *x11Native) open(win *Window, title string, width, height int) error {
	d, err := x11.Open(win, backend.Options{Title: title, Width: width, Height: height})
	if err != nil {
		return err
	}
	n.d = d
	return nil
}

func (n *x11Native) close() {
	if n.d != nil {
		n.d.Close()
	}
}

func (n *x11Native) pump(win *Window) { n.d.Pump(win) }
func (n *x11Native) present(win *Window, dirty canvas.Area) {
	n.d.Present(win, dirty)
}

// setFrameDeadline tells the driver when the next frame is due, so an idle
// Pump can sleep on the socket until then instead of spinning on it.
func (n *x11Native) setFrameDeadline(due time.Time) { n.d.SetFrameDeadline(due) }
func (n *x11Native) setTitle(title string)          { n.d.SetTitle(title) }
func (n *x11Native) setFullscreen(on bool) bool {
	return n.d.SetFullscreen(on)
}
func (n *x11Native) setOpacity(alpha uint8) bool          { return n.d.SetOpacity(alpha) }
func (n *x11Native) displaySize() (w, h int, ok bool)     { return n.d.DisplaySize() }
func (n *x11Native) displayRefresh() int                  { return n.d.DisplayRefresh() }
func (n *x11Native) setLimits(win *Window, limits Limits) { n.d.SetLimits(limits) }
func (n *x11Native) setSize(win *Window, width, height int) bool {
	return n.d.SetSize(win, width, height)
}
func (n *x11Native) displayScale() (float64, bool) { return n.d.DisplayScale() }

// systemDark takes the answer the desktop gives, which on X11 is the only
// answer there is: the server has no signal for the theme a desktop is
// drawing in, so it is read the way every toolkit reads it — out of the
// session's environment, the files the desktop writes at login and, where a
// desktop keeps it behind one, its own settings tool (see the linux
// package's SystemDark). A desktop with nothing to say still reports
// nothing, and [Window.SetSystemDark] stands in for it.
func (n *x11Native) systemDark() (bool, bool) { return x11.SystemDark() }

func (n *x11Native) clipboard() (text string, files []string) {
	return n.d.Clipboard()
}
func (n *x11Native) setIcon(images []*canvas.Canvas) bool { return n.d.SetIcon(images) }
func (n *x11Native) setClipboard(text string) bool        { return n.d.SetClipboard(text) }
