//go:build (linux && !android) || freebsd || openbsd || netbsd || dragonfly

package antui

func (n *x11Native) setCursor(icon CursorIcon) bool { return n.d.SetCursor(icon) }
