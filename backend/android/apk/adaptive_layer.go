package apk

import (
	_ "image/jpeg"

	"github.com/gabrielluizsf/antui/canvas"
)

// adaptiveLayer makes the foreground of an adaptive icon: the picture at two
// thirds of the layer, centred, with nothing around it.
func adaptiveLayer(src *canvas.Canvas, size int) *canvas.Canvas {
	out, err := canvas.NewCanvas(size, size)
	if err != nil {
		return nil
	}
	clear(out)
	inner := int(float64(size) * safeFraction)
	scaled := canvas.IconScaled(src, inner)
	if scaled == nil {
		return out
	}
	at := (size - inner) / 2
	out.Blit(at, at, scaled)
	return out
}
