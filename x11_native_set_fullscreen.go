//go:build (linux && !android) || freebsd || openbsd || netbsd || dragonfly

package antui

func (n *x11Native) setFullscreen(on bool) bool {
	return n.d.SetFullscreen(on)
}
