package svg

import (
	"math"

	"github.com/gabrielluizsf/antui/canvas"
)

// tileRegion is where the picture is drawn into, and at what scale: the whole
// tile, unless the whole tile at this scale is dearer than [maxTilePixels] is
// worth — then the part of the tile the canvas corners fall in, where even that
// fits, drawn at the same scale the pixels are sampled at; and where it does
// not, the whole tile at a scale that does, which is what every tile over the
// limit used to come out as. A pattern that draws past its own edges (an
// overflow=visible one) is tied to the tile's edges rather than to where the
// canvas looks, so it always takes the last road.
func tileRegion(x, y, w, h, s float64, visible bool, width, height int,
	inv, toTile canvas.Matrix) (rx, ry, rw, rh, scaled float64) {
	rx, ry, rw, rh, scaled = x, y, w, h, s
	if w*s*h*s <= maxTilePixels {
		return rx, ry, rw, rh, scaled
	}
	if !visible {
		// The box the four corners of the canvas fall in, taken back through
		// the same two maps every point goes through: everything outside it is
		// never asked about. A side the box runs past the tile on is the whole
		// of the tile that way — a box straddling the tile's edge wraps round
		// onto both sides of it.
		u0, v0, u1, v1 := math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)
		for _, c := range [][2]float64{
			{0, 0}, {float64(width), 0}, {0, float64(height)}, {float64(width), float64(height)},
		} {
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
	if rw*scaled*rh*scaled > maxTilePixels {
		// Even the part the canvas reaches is over the limit, or the pattern
		// draws past its own edges: the whole tile smaller.
		scaled *= math.Sqrt(maxTilePixels / (w * scaled * h * scaled))
		rx, ry, rw, rh = x, y, w, h
	}
	return rx, ry, rw, rh, scaled
}
