//go:build (linux && !android) || freebsd || openbsd || netbsd || dragonfly

package antui

func (n *x11Native) setOpacity(alpha uint8) bool { return n.d.SetOpacity(alpha) }
