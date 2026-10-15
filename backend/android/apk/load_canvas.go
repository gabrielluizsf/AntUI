package apk

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"os"

	"github.com/gabrielluizsf/antui/canvas"
)

// loadCanvas reads a PNG or a JPEG into a antui canvas.
func loadCanvas(path string) (*canvas.Canvas, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("apk: cannot read the icon: %w", err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("apk: %s is not a PNG or a JPEG: %w", path, err)
	}
	b := img.Bounds()
	cv, err := canvas.NewCanvas(b.Dx(), b.Dy())
	if err != nil {
		return nil, err
	}
	for y := range b.Dy() {
		for x := range b.Dx() {
			r, g, bl, a := img.At(b.Min.X+x, b.Min.Y+y).RGBA()
			// RGBA gives alpha-premultiplied 16-bit channels; a Canvas holds
			// straight 8-bit ones, so the multiplication has to be undone or
			// everything transparent comes out too dark.
			if a == 0 {
				cv.Put(x, y, canvas.Transparent)
				continue
			}
			cv.Put(x, y, canvas.RGBA(canvas.RGBAParams{R: uint8(r * 0xFF / a), G: uint8(g * 0xFF / a), B: uint8(bl * 0xFF / a), A: uint8(a >> 8)}))
		}
	}
	return cv, nil
}
