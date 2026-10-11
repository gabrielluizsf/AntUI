package svg

import (
	"math"

	"github.com/gabrielluizsf/antui/canvas"
)

// patternShade is the colour one pixel of a pattern-painted shape takes: the
// tile is drawn into a picture of its own, and the pixel asks that picture
// which part of the tile it is over — and the same part is asked for every
// pixel all the way across the shape, wrapped round by the size of the tile.
//
// The pixel is taken back to the coordinates the shape was written in, the same
// as for a gradient, and then backwards through `patternTransform` — because
// that transform moves the tile rather than the shape, so where the tile now
// stands is what has to be found. What comes out is a place inside the tile,
// worked out before the transform moved it, which is where the picture was
// drawn. A transform with no way back, a tile of no size and a pattern that is
// already being drawn further out all answer the same way: there is no shade
// here, and the shape takes the colour after the `url(...)` if it brought one
// with it and nothing at all if it did not.
func patternShade(pt *pattern, n *Node, m canvas.Matrix, width, height int, current canvas.Color, opacity float64, patterning []*pattern) (func(x, y int) canvas.Color, bool) {
	for _, being := range patterning {
		if being == pt {
			// This pattern is already being drawn further up: following it
			// again would draw it while it is being drawn, for ever — a
			// pattern that names itself, or two that name each other, is left
			// off so the drawing comes out as a picture rather than not at all.
			return nil, false
		}
	}
	box := n.box
	x, y, w, h := pt.x, pt.y, pt.w, pt.h
	if pt.boxUnits {
		if box[2] <= 0 || box[3] <= 0 {
			return nil, false
		}
		x, y = box[0]+x*box[2], box[1]+y*box[3]
		w, h = w*box[2], h*box[3]
	}
	if pt.contentBox && (box[2] <= 0 || box[3] <= 0) {
		return nil, false
	}
	if w <= 0 || h <= 0 {
		return nil, false
	}
	toCanvas := m.Mul(n.Style.Transform)
	inv, ok := toCanvas.Inverse()
	if !ok {
		return nil, false
	}
	toTile, ok := pt.transform.Inverse()
	if !ok {
		return nil, false
	}
	// How much of the canvas one step of the tile is worth, which is the size
	// the tile lands at: the drawing's own scale, put through the transform
	// that moves the pattern, since that transform moves the tile too. A tile
	// worth more than [maxTilePixels] comes out drawn as the part of it the
	// canvas can ask about at this same scale, and smaller only where even
	// that is over the limit.
	s := scaleOf(toCanvas.Mul(pt.transform))
	if !(s > 0) || math.IsInf(s, 0) {
		return nil, false
	}
	// Where the picture is drawn into: the whole tile, unless the whole tile
	// at this scale is dearer than [maxTilePixels] is worth — then the part of
	// the tile the canvas corners fall in, where even that fits, drawn at the
	// same scale the pixels are sampled at; and where it does not, the whole
	// tile at a scale that does, which is what every tile over the limit used
	// to come out as. A pattern that draws past its own edges (an
	// overflow=visible one) is tied to the tile's edges rather than to where
	// the canvas looks, so it always takes the last road.
	rx, ry, rw, rh := x, y, w, h
	if w*s*h*s > maxTilePixels {
		if !pt.visible {
			// The box the four corners of the canvas fall in, taken back
			// through the same two maps every point goes through: everything
			// outside it is never asked about. A side the box runs past the
			// tile on is the whole of the tile that way — a box straddling
			// the tile's edge wraps round onto both sides of it.
			u0, v0, u1, v1 := math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)
			for _, c := range [][2]float64{{0, 0}, {float64(width), 0}, {0, float64(height)}, {float64(width), float64(height)}} {
				ux, uy := inv.Map(c[0], c[1])
				bx, by := toTile.Map(ux, uy)
				u0, v0 = min(u0, bx-x), min(v0, by-y)
				u1, v1 = max(u1, bx-x), max(v1, by-y)
			}
			if u0 >= 0 && u1 <= w {
				rx, rw = x+u0, u1-u0
			}
			if v0 >= 0 && v1 <= h {
				ry, rh = y+v0, v1-v0
			}
		}
		if rw*s*rh*s > maxTilePixels {
			// Even the part the canvas reaches is over the limit, or the
			// pattern draws past its own edges: the whole tile smaller.
			s *= math.Sqrt(maxTilePixels / (w * s * h * s))
			rx, ry, rw, rh = x, y, w, h
		}
	}
	// The picture is cut to the part above, so a pixel takes its colour from
	// the tile at the same scale minus how far into the tile the part starts.
	xoff, yoff := (rx-x)*s, (ry-y)*s
	tile := pt.picture(box, rx, ry, rw, rh, s, current, patterning)
	if tile == nil {
		return nil, false
	}
	return func(px, py int) canvas.Color {
		ux, uy := inv.Map(float64(px)+0.5, float64(py)+0.5)
		bx, by := toTile.Map(ux, uy)
		ix := int(modUnit(bx-x, w)*s - xoff)
		iy := int(modUnit(by-y, h)*s - yoff)
		if ix < 0 {
			ix = 0
		} else if ix >= tile.Width {
			ix = tile.Width - 1
		}
		if iy < 0 {
			iy = 0
		} else if iy >= tile.Height {
			iy = tile.Height - 1
		}
		return fade(tile.At(ix, iy), opacity)
	}, true
}
