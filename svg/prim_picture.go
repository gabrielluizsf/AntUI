package svg

import (
	"math"

	"github.com/gabrielluizsf/antui/canvas"
)

// primPicture is what one primitive leaves behind: its own picture of the
// pictures it took, never one of them itself, because the same picture may be
// named again by what follows and a primitive must not change a picture it
// does not own. The lengths a primitive writes for itself — an offset, a
// blur's radius — are turned into pixels here, the way a blur's in the
// attribute is turned in [applyFilters]: one unit of the drawing is scale
// pixels of the canvas, per axis where the length runs along one.
func primPicture(p filterPrim, in, in2, src *canvas.Canvas, s scale2, current canvas.Color) *canvas.Canvas {
	switch p.kind {
	case primGaussianBlur:
		out := pictureOf(in)
		if out == nil {
			return nil
		}
		out.BlurXY(canvas.Area{X: 0, Y: 0, Width: out.Width, Height: out.Height}, radius(p.stdDev*s.x), radius(p.stdDevY*s.y))
		return out
	case primColorMatrix:
		out := pictureOf(in)
		colorMatrix(out, p.matrix)
		return out
	case primOffset:
		out := blankPicture(src.Width, src.Height)
		if out == nil || in == nil {
			return out
		}
		out.BlitOver(int(math.Round(p.dx*s.x)), int(math.Round(p.dy*s.y)), in)
		return out
	case primFlood:
		out := blankPicture(src.Width, src.Height)
		if out == nil {
			return nil
		}
		c := p.flood
		if p.floodCurrent {
			c = current
		}
		out.FillRect(canvas.Area{X: 0, Y: 0, Width: out.Width, Height: out.Height}, fade(c, p.floodOpacity))
		return out
	case primComposite:
		out := pictureOf(in)
		pairMap(out, in2, func(a, b canvas.Color) canvas.Color { return compositePx(a, b, p.comp, p.k) })
		return out
	case primBlend:
		out := pictureOf(in)
		pairMap(out, in2, func(a, b canvas.Color) canvas.Color { return blendPx(a, b, p.blend) })
		return out
	}
	return nil
}
