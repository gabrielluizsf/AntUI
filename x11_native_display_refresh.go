//go:build (linux && !android) || freebsd || openbsd || netbsd || dragonfly

package antui

func (n *x11Native) displayRefresh() int { return n.d.DisplayRefresh() }
