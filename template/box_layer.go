package template

import "github.com/gabrielluizsf/antui/canvas"

// boxLayer makes the empty canvas a background layer is rasterized into. A
// window-size canvas is the rare allocation that can fail; a nil result makes
// the caller skip the layer.
func boxLayer(w, h int) *canvas.Canvas {
	cv, err := canvas.NewCanvas(w, h)
	if err != nil {
		return nil
	}
	return cv
}
