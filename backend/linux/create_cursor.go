//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"encoding/binary"

	"github.com/gabrielluizsf/antui/backend"
	"github.com/gabrielluizsf/antui/canvas"
)

// createCursor builds an X cursor from icon's picture and returns its id.
// X cursors are two colours: a pixel is foreground where the picture is
// bright and opaque, background where it is dark and opaque, and invisible
// everywhere else — with the invisible pixels bordering an opaque one kept
// background too, a one-pixel dark edge around the shape, which is what
// keeps a white glyph readable over a white page instead of vanishing into
// it. The foreground itself is the average of the picture's own bright
// pixels rather than a fixed white, so a neon-blue icon stays neon blue
// instead of bleaching to the one colour every legacy cursor defaults to.
func (x *Driver) createCursor(icon backend.CursorIcon) uint32 {
	img := icon.Image
	w, h := img.Width, img.Height

	visible := func(px, py int) bool {
		if px < 0 || py < 0 || px >= w || py >= h {
			return false
		}
		return img.At(px, py).A() > 127
	}
	bright := func(px, py int) bool {
		if !visible(px, py) {
			return false
		}
		c := img.At(px, py)
		return int(c.R())+int(c.G())+int(c.B()) >= 384 // half of 3*255
	}
	outline := func(px, py int) bool {
		if visible(px, py) {
			return true
		}
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				if visible(px+dx, py+dy) {
					return true
				}
			}
		}
		return false
	}

	foreR, foreG, foreB := averageBright(img, w, h, bright)

	source := x.createPixmap(1, w, h)
	mask := x.createPixmap(1, w, h)
	gc := x.createGCFor(source)
	x.putBitmap(source, gc, x.packBitmap(bright, w, h), w, h)
	x.putBitmap(mask, gc, x.packBitmap(outline, w, h), w, h)
	x.freeGC(gc)

	cursor := x.newID()
	body := make([]byte, 28)
	binary.LittleEndian.PutUint32(body[0:], cursor)
	binary.LittleEndian.PutUint32(body[4:], source)
	binary.LittleEndian.PutUint32(body[8:], mask)
	binary.LittleEndian.PutUint16(body[12:], uint16(foreR)<<8|uint16(foreR)) // fore red
	binary.LittleEndian.PutUint16(body[14:], uint16(foreG)<<8|uint16(foreG)) // fore green
	binary.LittleEndian.PutUint16(body[16:], uint16(foreB)<<8|uint16(foreB)) // fore blue
	binary.LittleEndian.PutUint16(body[18:], 0)                              // back red
	binary.LittleEndian.PutUint16(body[20:], 0)                              // back green
	binary.LittleEndian.PutUint16(body[22:], 0)                              // back blue
	binary.LittleEndian.PutUint16(body[24:], uint16(icon.HotX))
	binary.LittleEndian.PutUint16(body[26:], uint16(icon.HotY))
	x.send(xCreateCursor, 0, body)

	x.freePixmap(source)
	x.freePixmap(mask)
	return cursor
}

// averageBright is the mean colour of every pixel bright reports true for —
// the one foreground colour a legacy cursor gets to paint all of them in.
// White, the shade every one of those pixels was picked for being closest
// to, when the picture has none.
func averageBright(img *canvas.Canvas, w, h int, bright func(px, py int) bool) (r, g, b uint8) {
	var sumR, sumG, sumB, count int
	for py := 0; py < h; py++ {
		for px := 0; px < w; px++ {
			if !bright(px, py) {
				continue
			}
			c := img.At(px, py)
			sumR += int(c.R())
			sumG += int(c.G())
			sumB += int(c.B())
			count++
		}
	}
	if count == 0 {
		return 0xFF, 0xFF, 0xFF
	}
	return uint8(sumR / count), uint8(sumG / count), uint8(sumB / count)
}
