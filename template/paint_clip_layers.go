package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// paintClipLayers paints a clipped box's background images: each layer is
// rasterized as usual, then masked to the clip-path polygon, because a canvas
// clip is rectangular and the site's chamfer is not.
func (cs *cssStyle) paintClipLayers(win *antui.Window, st css.Style, x, y, w, h int, pts []canvas.Point) {
	cv := win.Canvas()
	for i := len(st.BackgroundImages) - 1; i >= 0; i-- {
		layer := cs.paintBackgroundLayer(win, st, i, x, y, w, h, 0)
		if layer == nil {
			continue
		}
		maskClip(layer, pts, x, y)
		cv.BlitOver(x, y, layer)
	}
}

// maskClip zeroes every pixel of a layer outside the polygon, so blitting it
// paints only the clipped shape.
func maskClip(layer *canvas.Canvas, pts []canvas.Point, x, y int) {
	if layer.Pixels == nil {
		return
	}
	for py := range layer.Height {
		for px := range layer.Width {
			if !insidePoly(pts, float64(x+px)+0.5, float64(y+py)+0.5) {
				layer.Pixels[py*layer.Stride+px] = 0
			}
		}
	}
}

// insidePoly reports whether a point is inside a polygon by the even-odd rule.
func insidePoly(pts []canvas.Point, x, y float64) bool {
	inside := false
	for i, j := 0, len(pts)-1; i < len(pts); j, i = i, i+1 {
		yi, yj, xi, xj := pts[i].Y, pts[j].Y, pts[i].X, pts[j].X
		if (yi > y) != (yj > y) && x < (xj-xi)*(y-yi)/(yj-yi)+xi {
			inside = !inside
		}
	}
	return inside
}
