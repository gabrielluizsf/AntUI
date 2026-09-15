//go:build android

package antui

import (
	"github.com/gabrielluizsf/antui/backend/android/ndk"
	"github.com/gabrielluizsf/antui/canvas"
)

// packRGBA swizzles a canvas color into the 0xAABBGGRR word Android's
// RGBA_8888 format expects on a little-endian machine.
func packRGBA(c canvas.Color) uint32 {
	v := uint32(c)
	return v&0xFF00FF00 | (v&0x00FF0000)>>16 | (v&0x000000FF)<<16
}

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
