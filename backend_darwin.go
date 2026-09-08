//go:build darwin

package antui

import (
	"github.com/gabrielluizsf/antui/backend"
	darwin "github.com/gabrielluizsf/antui/backend/macos"
	"github.com/gabrielluizsf/antui/canvas"
)

// darwinNative is the macOS backend, in two halves: this struct adapts the
// window core's calls to the platform package, and the package itself lives
// in backend/macos.
type darwinNative struct{ d *darwin.Driver }

// newBackend is what OpenWith asks for when there is not yet a window.
func newBackend() platform { return &darwinNative{} }

func (n *darwinNative) open(win *Window, title string, width, height int) error {
	d, err := darwin.Open(win, backend.Options{Title: title, Width: width, Height: height})
	if err != nil {
		return err
	}
	n.d = d
	return nil
}

func (n *darwinNative) close() {
	if n.d != nil {
		n.d.Close()
	}
}

func (n *darwinNative) pump(win *Window) { n.d.Pump(win) }
func (n *darwinNative) present(win *Window, dirty canvas.Area) {
	n.d.Present(win, dirty)
}
func (n *darwinNative) setTitle(title string) { n.d.SetTitle(title) }
func (n *darwinNative) setFullscreen(on bool) bool {
	return n.d.SetFullscreen(on)
}
func (n *darwinNative) displaySize() (w, h int, ok bool) { return n.d.DisplaySize() }
func (n *darwinNative) displayRefresh() int              { return n.d.DisplayRefresh() }
func (n *darwinNative) setLimits(win *Window, limits Limits) { n.d.SetLimits(limits) }
func (n *darwinNative) setSize(win *Window, width, height int) bool {
	return n.d.SetSize(win, width, height)
}
func (n *darwinNative) contentScale() float64 { return n.d.ContentScale() }
func (n *darwinNative) clipboard() (text string, files []string) {
	return n.d.Clipboard()
}
func (n *darwinNative) setIcon(images []*canvas.Canvas) bool { return n.d.SetIcon(images) }
func (n *darwinNative) setClipboard(text string) bool        { return n.d.SetClipboard(text) }