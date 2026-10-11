//go:build (linux && !android) || freebsd || openbsd || netbsd || dragonfly

package antui

func (n *x11Native) setSize(win *Window, width, height int) bool {
	return n.d.SetSize(win, width, height)
}
