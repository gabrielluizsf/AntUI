package template

import "github.com/gabrielluizsf/antui/template/css"

// gridPlaceOut hands an out-of-flow item the box solve worked out for it.
func (c *CSS) gridPlaceOut(st css.Style, it gridItem) (int, int, int, int, css.Style, bool) {
	return c.grid.originX + it.x, c.grid.originY + it.y, it.w, it.h, st, true
}
