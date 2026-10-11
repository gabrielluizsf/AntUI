//go:build (linux && !android) || freebsd || openbsd || netbsd || dragonfly

package antui

func (n *x11Native) setTitle(title string) { n.d.SetTitle(title) }
