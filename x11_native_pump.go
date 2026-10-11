//go:build (linux && !android) || freebsd || openbsd || netbsd || dragonfly

package antui

func (n *x11Native) pump(win *Window) { n.d.Pump(win) }
