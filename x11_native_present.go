//go:build (linux && !android) || freebsd || openbsd || netbsd || dragonfly

package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

func (n *x11Native) present(win *Window, dirty canvas.Area) bool {
	return n.d.Present(win, dirty)
}
