package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// applyFilters runs a filter list over a rectangle of the window. A blur
// length is scaled with the window; drop-shadow is left to paintShadows.
func (cs *cssStyle) applyFilters(cv *canvas.Canvas, filters []css.Filter, x, y, w, h int) {
	u := cs.u()
	for _, f := range filters {
		if f.Drop != nil {
			continue
		}
		amount := f.Amount
		if f.Kind == canvas.FilterBlur {
			amount = f.Amount * float64(u)
		}
		cv.FilterRegion(x, y, w, h, f.Kind, amount)
	}
}
