//go:build android

package antui

import (
	"github.com/gabrielluizsf/antui/backend/android/ndk"
	"github.com/gabrielluizsf/antui/canvas"
)

// blit8888 is the path every device takes. RGBX is the same bytes as RGBA
// with the alpha ignored, so it is the same swizzle — writing the alpha into
// a channel the compositor does not read costs nothing and keeps one loop.
func blit8888(buf ndk.Buffer, cv *canvas.Canvas, x0, y0, x1, y1 int) {
	dst := buf.Pixels32()
	if dst == nil {
		return
	}
	for y := y0; y < y1; y++ {
		d := dst[y*buf.Stride+x0 : y*buf.Stride+x1]
		s := cv.Pixels[y*cv.Stride+x0 : y*cv.Stride+x1]
		d = d[:len(s)]
		for i, c := range s {
			d[i] = packRGBA(c)
		}
	}
}
