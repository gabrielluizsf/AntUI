package svg

import "github.com/gabrielluizsf/antui/canvas"

// tileMap takes a canvas pixel back to a place inside the tile: through the
// drawing to the shape's coordinates, and then backwards through the pattern's
// transform, which is the map that moves the tile rather than the shape.
type tileMap struct {
	inv, toTile   canvas.Matrix
	tile          *canvas.Canvas
	x, y, w, h    float64
	s, xoff, yoff float64
	opacity       float64
}

// at is the colour one pixel of the shape takes, wrapped round by the size of
// the tile and clamped to the picture where it falls over its edge.
func (tm tileMap) at(px, py int) canvas.Color {
	ux, uy := tm.inv.Map(float64(px)+0.5, float64(py)+0.5)
	bx, by := tm.toTile.Map(ux, uy)
	ix := min(max(int(modUnit(bx-tm.x, tm.w)*tm.s-tm.xoff), 0), tm.tile.Width-1)
	iy := min(max(int(modUnit(by-tm.y, tm.h)*tm.s-tm.yoff), 0), tm.tile.Height-1)
	return fade(tm.tile.At(ix, iy), tm.opacity)
}

// tileMapFor maps a picture drawn for the region (rx, ry) at scale s, inside a
// tile at (x, y) of (w, h).
func tileMapFor(inv, toTile canvas.Matrix, tile *canvas.Canvas,
	x, y, w, h, s, rx, ry, opacity float64) tileMap {
	return tileMap{
		inv: inv, toTile: toTile, tile: tile,
		x: x, y: y, w: w, h: h, s: s,
		xoff: (rx - x) * s, yoff: (ry - y) * s, opacity: opacity,
	}
}
