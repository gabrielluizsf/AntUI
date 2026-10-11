package antui

import (
	"bytes"
	"github.com/gabrielluizsf/antui/canvas"
	"image"
	"image/png"
)

// encodeIconPNG writes a picture as a PNG, which is what an icon over 48
// pixels is these days and what every Windows since Vista reads. The old
// format is a bitmap with the rows upside down and a mask nobody uses; there
// is no reason to write one.
func encodeIconPNG(image_ *canvas.Canvas) ([]byte, error) {
	out := image.NewNRGBA(image.Rect(0, 0, image_.Width, image_.Height))
	for y := range image_.Height {
		for x := range image_.Width {
			c := image_.At(x, y)
			at := out.PixOffset(x, y)
			out.Pix[at+0] = c.R()
			out.Pix[at+1] = c.G()
			out.Pix[at+2] = c.B()
			out.Pix[at+3] = c.A()
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, out); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
