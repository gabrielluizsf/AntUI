package svg

import "github.com/gabrielluizsf/antui/canvas"

// drawTiledPicture lays the copies of a picture whose overflow reaches into
// the tile, painted left to right in rows from the top — the order the spec
// says tiles are laid in, which is what says which of two overlapping copies is
// on top. The picture's span in the tile's own coordinates comes out of what it
// painted, through the same box the picture itself went through. It reports
// whether any copy reaches in.
func (pt *pattern) drawTiledPicture(box [4]float64, x, y, w, h float64, base, boxT canvas.Matrix, one func(canvas.Matrix)) bool {
	b, ok := paintedBox(pt.content)
	if !ok {
		return false
	}
	x0, y0, x1, y1 := b[0], b[1], b[0]+b[2], b[1]+b[3]
	if pt.contentBox {
		x0, x1 = box[0]+x0*box[2], box[0]+x1*box[2]
		y0, y1 = box[1]+y0*box[3], box[1]+y1*box[3]
	}
	nx0, nx1, ok := copyRange(x0-x, x1-x, w)
	if !ok {
		return false
	}
	ny0, ny1, ok := copyRange(y0-y, y1-y, h)
	if !ok {
		return false
	}
	for ny := ny0; ny <= ny1; ny++ {
		for nx := nx0; nx <= nx1; nx++ {
			one(base.Mul(canvas.Translate(float64(nx)*w, float64(ny)*h)).Mul(boxT))
		}
	}
	return true
}
