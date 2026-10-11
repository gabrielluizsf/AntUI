//go:build android

package antui

import (
	"github.com/gabrielluizsf/antui/backend/android/ndk"
	"github.com/gabrielluizsf/antui/canvas"
)

// blit copies a rectangle of the canvas into a locked window buffer.
//
// Two things about that buffer are easy to get wrong and cost a day each.
//
// Its stride is in **pixels, not bytes**, and it is not the width: the
// compositor rounds rows up, often to a multiple of 16 or 32. Using the
// width to step between rows shears the picture a little further every row,
// which looks like a diagonal tear and reads like a memory bug.
//
// Its format is named in **byte order**. Android's RGBA_8888 is R, G, B, A
// in memory, so on a little-endian machine the word is 0xAABBGGRR — while a
// [Color] is 0xAARRGGBB. Red and blue are the wrong way round, and a picture
// that comes out looking like a photographic negative of itself is this and
// nothing else.
func blit(buf ndk.Buffer, cv *canvas.Canvas, r ndk.Rect) {
	if cv == nil || cv.Pixels == nil {
		// The window's canvas is always the wide format, so this cannot
		// happen from inside the library. It is here so that a narrow canvas
		// arriving from somewhere new leaves a black frame rather than a
		// crash on a phone in someone's hand.
		return
	}
	x0 := int(max(r.Left, 0))
	y0 := int(max(r.Top, 0))
	x1 := min(int(r.Right), min(buf.Width, cv.Width))
	y1 := min(int(r.Bottom), min(buf.Height, cv.Height))
	if x0 >= x1 || y0 >= y1 {
		return
	}

	switch buf.Format {
	case ndk.RGBA8888, ndk.RGBX8888:
		blit8888(buf, cv, x0, y0, x1, y1)
	case ndk.RGB565:
		blit565(buf, cv, x0, y0, x1, y1)
	default:
		// RGB_888 exists and no device has shipped with it in a decade. A
		// black frame and a line in the log beats a wrong-looking one nobody
		// can explain.
		ndk.Warnf("no blit for surface format %s", buf.Format)
	}
}
