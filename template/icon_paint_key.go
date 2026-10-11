package template

import "github.com/gabrielluizsf/antui/canvas"

// iconPaintKey is one picture: which icon, how big, in which colour.
type iconPaintKey struct {
	name string
	w, h int
	tint canvas.Color
}
