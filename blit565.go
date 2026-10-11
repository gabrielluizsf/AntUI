//go:build android

package antui

import (
	"github.com/gabrielluizsf/antui/backend/android/ndk"
	"github.com/gabrielluizsf/antui/canvas"
)

// blit565 goes through canvas.To565, the same packer the Canvas's own 565 format
// uses, so a colour reaches a 16-bit surface as the same word whichever way
// it got there. It rounds rather than truncating, which is what keeps white
// white.
//
// It does not dither. The Canvas has a dithered path for its own narrow
// format; a *surface* that is 565 belongs to a device old enough that a
// second pass over every pixel is the thing that would be noticed.
func blit565(buf ndk.Buffer, cv *canvas.Canvas, x0, y0, x1, y1 int) {
	dst := buf.Pixels16()
	if dst == nil {
		return
	}
	for y := y0; y < y1; y++ {
		d := dst[y*buf.Stride+x0 : y*buf.Stride+x1]
		s := cv.Pixels[y*cv.Stride+x0 : y*cv.Stride+x1]
		d = d[:len(s)]
		for i, c := range s {
			d[i] = canvas.To565(c)
		}
	}
}
