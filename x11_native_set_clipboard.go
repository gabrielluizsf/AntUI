//go:build (linux && !android) || freebsd || openbsd || netbsd || dragonfly

package antui

func (n *x11Native) setClipboard(text string) bool { return n.d.SetClipboard(text) }
