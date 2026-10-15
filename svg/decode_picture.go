package svg

import (
	"bytes"
	"fmt"
	"image"

	"github.com/gabrielluizsf/antui/canvas"
)

// decodePicture is a PNG, a JPEG or a GIF turned into a canvas, at the size it
// was encoded at. The channels come out of the decoder premultiplied by their
// alpha and a canvas holds them straight, so the multiplication is undone the
// same way the icon reader in the android package does it: everything
// transparent would otherwise come out too dark to lay over anything.
func decodePicture(data []byte) (*canvas.Canvas, error) {
	src, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("it is not a PNG, a JPEG or a GIF: %w", err)
	}
	b := src.Bounds()
	if b.Dx() <= 0 || b.Dy() <= 0 {
		return nil, fmt.Errorf("it has no size to it")
	}
	cv, err := canvas.NewCanvas(b.Dx(), b.Dy())
	if err != nil {
		return nil, err
	}
	for y := range b.Dy() {
		for x := range b.Dx() {
			r, g, bl, a := src.At(b.Min.X+x, b.Min.Y+y).RGBA()
			if a == 0 {
				cv.Put(x, y, canvas.Transparent)
				continue
			}
			cv.Put(x, y, canvas.RGBA(canvas.RGBAParams{R: uint8(r * 0xFF / a), G: uint8(g * 0xFF / a), B: uint8(bl * 0xFF / a), A: uint8(a >> 8)}))
		}
	}
	return cv, nil
}
