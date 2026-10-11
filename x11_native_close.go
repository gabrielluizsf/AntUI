//go:build (linux && !android) || freebsd || openbsd || netbsd || dragonfly

package antui

func (n *x11Native) close() {
	if n.d != nil {
		n.d.Close()
	}
}
