//go:build (linux && !android) || freebsd || openbsd || netbsd || dragonfly

package antui

func (n *x11Native) displaySize() (w, h int, ok bool) { return n.d.DisplaySize() }
