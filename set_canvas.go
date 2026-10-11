package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// SetCanvas swaps the window's pixel surface for cv and returns the previous
// one, so a caller can paint into a scratch canvas and hand it back. The
// canvas must have the same size as the window — the width and height the
// window reports and the pixels it redraws stay pinned to that size — while
// the surface in between may hold a partially composed scene, like a layer
// waiting to be blitted over.
func (win *Window) SetCanvas(cv *canvas.Canvas) *canvas.Canvas {
	if cv == nil {
		return win.cv
	}
	prev := win.cv
	win.cv = cv
	return prev
}
