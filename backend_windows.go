//go:build windows

package antui

import (
	"github.com/gabrielluizsf/antui/backend"
	win32 "github.com/gabrielluizsf/antui/backend/windows"
	"github.com/gabrielluizsf/antui/canvas"
)

// windowsNative is the Win32 backend, in two halves: this struct adapts the
// window core's calls to the platform package, and the package itself lives
// in backend/windows.
type windowsNative struct{ d *win32.Driver }

// newBackend is what OpenWith asks for when there is not yet a window.
func newBackend() platform { return &windowsNative{} }

func (n *windowsNative) open(win *Window, title string, width, height int) error {
	d, err := win32.Open(win, backend.Options{Title: title, Width: width, Height: height})
	if err != nil {
		return err
	}
	n.d = d
	return nil
}

func (n *windowsNative) close() {
	if n.d != nil {
		n.d.Close()
	}
}

func (n *windowsNative) pump(win *Window) { n.d.Pump(win) }
func (n *windowsNative) present(win *Window, dirty canvas.Area) {
	n.d.Present(win, dirty)
}
func (n *windowsNative) setTitle(title string) { n.d.SetTitle(title) }
func (n *windowsNative) setFullscreen(on bool) bool {
	return n.d.SetFullscreen(on)
}
func (n *windowsNative) displaySize() (w, h int, ok bool) { return n.d.DisplaySize() }
func (n *windowsNative) displayRefresh() int              { return n.d.DisplayRefresh() }
func (n *windowsNative) setLimits(win *Window, limits Limits) { n.d.SetLimits(limits) }
func (n *windowsNative) setSize(win *Window, width, height int) bool {
	return n.d.SetSize(win, width, height)
}
func (n *windowsNative) contentScale() float64 { return n.d.ContentScale() }
func (n *windowsNative) clipboard() (text string, files []string) {
	return n.d.Clipboard()
}
func (n *windowsNative) setIcon(images []*canvas.Canvas) bool { return n.d.SetIcon(images) }
func (n *windowsNative) setClipboard(text string) bool        { return n.d.SetClipboard(text) }