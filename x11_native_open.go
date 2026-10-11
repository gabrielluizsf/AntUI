//go:build (linux && !android) || freebsd || openbsd || netbsd || dragonfly

package antui

import (
	"github.com/gabrielluizsf/antui/backend"
	x11 "github.com/gabrielluizsf/antui/backend/linux"
)

func (n *x11Native) open(win *Window, title string, width, height int) error {
	d, err := x11.Open(win, backend.Options{Title: title, Width: width, Height: height})
	if err != nil {
		return err
	}
	n.d = d
	return nil
}
