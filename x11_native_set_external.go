//go:build (linux && !android) || freebsd || openbsd || netbsd || dragonfly

package antui

func (n *x11Native) setExternal(on bool) bool { return n.d.SetExternal(on) }
