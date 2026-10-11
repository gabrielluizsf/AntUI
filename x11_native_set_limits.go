//go:build (linux && !android) || freebsd || openbsd || netbsd || dragonfly

package antui

func (n *x11Native) setLimits(win *Window, limits Limits) { n.d.SetLimits(limits) }
