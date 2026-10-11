//go:build (linux && !android) || freebsd || openbsd || netbsd || dragonfly

package antui

func (n *x11Native) displayScale() (float64, bool) { return n.d.DisplayScale() }
