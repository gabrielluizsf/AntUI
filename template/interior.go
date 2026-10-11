package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// interior is the colour behind the border: the box's own fill when it has one,
// and nothing at all when it does not, so a sheet that gives a border and no
// surface leaves the page showing through the middle of the ring instead of the
// theme's colour painted over it.
func (cs *cssStyle) interior(st css.Style, hasBg bool) canvas.Color {
	if hasBg {
		return cs.fill(st)
	}
	return canvas.Transparent
}
