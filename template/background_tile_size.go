package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// backgroundTileSize scales the source to the layer: background-size chooses
// cover, contain, or the given and auto-adjusted dimensions.
func (cs *cssStyle) backgroundTileSize(st css.Style, i int, src *canvas.Canvas, ow, oh int) (int, int) {
	sz := propAt(st.BackgroundSize, i, css.BackSize{})
	dw, dh := src.Width, src.Height
	switch {
	case sz.Cover:
		dw, dh = coverSize(src.Width, src.Height, ow, oh)
	case sz.Contain:
		dw, dh = canvas.Fit(src.Width, src.Height, canvas.Area{Width: ow, Height: oh})
	default:
		if !autoDim(sz.W) {
			dw = cs.sizeDim(st, sz.W, ow, src.Width)
		}
		if !autoDim(sz.H) {
			dh = cs.sizeDim(st, sz.H, oh, src.Height)
		}
	}
	return dw, dh
}
