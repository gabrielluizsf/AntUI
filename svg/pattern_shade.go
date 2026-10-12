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
func patternShade(pt *pattern, n *Node, m canvas.Matrix, width, height int,
	current canvas.Color, opacity float64, patterning []*pattern) (func(x, y int) canvas.Color, bool) {
	if pt.beingDrawn(patterning) {
		return nil, false
	}
	x, y, w, h, ok := tileBox(pt, n)
	if !ok {
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
	// that moves the pattern, since that transform moves the tile too.
	s := scaleOf(toCanvas.Mul(pt.transform))
	if !(s > 0) || math.IsInf(s, 0) {
		return nil, false
	}
	rx, ry, rw, rh, s := tileRegion(x, y, w, h, s, pt.visible, width, height, inv, toTile)
	// The picture is cut to that part, so a pixel takes its colour from the
	// tile at the same scale minus how far into the tile the part starts.
	tile := pt.picture(n.box, rx, ry, rw, rh, s, current, patterning)
	if tile == nil {
		return nil, false
	}
	tm := tileMapFor(inv, toTile, tile, x, y, w, h, s, rx, ry, opacity)
	return tm.at, true
}
