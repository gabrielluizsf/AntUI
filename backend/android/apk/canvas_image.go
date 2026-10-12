package apk

import (
	_ "image/jpeg"

	"github.com/gabrielluizsf/antui/canvas"
)

// canvasImage is a canvas.Canvas seen as an image.Image, which is what the
// standard library's encoder wants. Canvas has At, but it answers a
// canvas.Color rather than a color.Color, so it cannot be the interface
// itself.
type canvasImage struct{ cv *canvas.Canvas }
