//go:build linux || freebsd || openbsd || netbsd || dragonfly

package antui

import (
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
func (n *x11Native) setTitle(title string) { n.d.SetTitle(title) }
func (n *x11Native) setFullscreen(on bool) bool {
	return n.d.SetFullscreen(on)
}
func (n *x11Native) displaySize() (w, h int, ok bool) { return n.d.DisplaySize() }
func (n *x11Native) displayRefresh() int              { return n.d.DisplayRefresh() }
func (n *x11Native) setLimits(win *Window, limits Limits) { n.d.SetLimits(limits) }
func (n *x11Native) setSize(win *Window, width, height int) bool {
	return n.d.SetSize(win, width, height)
}
func (n *x11Native) contentScale() float64 { return n.d.ContentScale() }
func (n *x11Native) clipboard() (text string, files []string) {
	return n.d.Clipboard()
}
func (n *x11Native) setIcon(images []*canvas.Canvas) bool { return n.d.SetIcon(images) }
func (n *x11Native) setClipboard(text string) bool        { return n.d.SetClipboard(text) }