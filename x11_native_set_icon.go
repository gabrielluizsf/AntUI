//go:build (linux && !android) || freebsd || openbsd || netbsd || dragonfly

package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

func (n *x11Native) setIcon(images []*canvas.Canvas) bool { return n.d.SetIcon(images) }
