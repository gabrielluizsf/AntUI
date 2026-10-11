//go:build (linux && !android) || freebsd || openbsd || netbsd || dragonfly

package antui

func (n *x11Native) clipboard() (text string, files []string) {
	return n.d.Clipboard()
}
