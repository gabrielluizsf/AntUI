package template

import "github.com/gabrielluizsf/antui/canvas"

// blitTiled paints the image over the layer, repeating the tile from the
// anchor across the whole layer.
func blitTiled(layer *canvas.Canvas, src *canvas.Canvas, px, py, dw, dh, w, h int) {
	for tx := boxStart(px, dw); tx < w; tx += dw {
		for ty := boxStart(py, dh); ty < h; ty += dh {
			layer.BlitScaled(tx, ty, dw, dh, src)
		}
	}
}
