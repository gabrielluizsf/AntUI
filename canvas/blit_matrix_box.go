package canvas

import "math"

// blitMatrixBox is the destination bbox of a blit: the footprint of the
// region's corners through m, clamped to the canvas, and whether any of it is
// inside the clip region.
func blitMatrixBox(cv *Canvas, m Matrix, x0, y0, x1, y1 int) (int, int, int, int, bool) {
	ax, ay := m.Map(float64(x0), float64(y0))
	bx, by := m.Map(float64(x1), float64(y0))
	cx, cy := m.Map(float64(x1), float64(y1))
	dx, dy := m.Map(float64(x0), float64(y1))
	minX := math.Min(math.Min(ax, bx), math.Min(cx, dx))
	maxX := math.Max(math.Max(ax, bx), math.Max(cx, dx))
	minY := math.Min(math.Min(ay, by), math.Min(cy, dy))
	maxY := math.Max(math.Max(ay, by), math.Max(cy, dy))

	d0x := max(int(math.Floor(minX)), cv.Clip.X)
	d1x := min(int(math.Ceil(maxX)), cv.Clip.X+cv.Clip.Width)
	d0y := max(int(math.Floor(minY)), cv.Clip.Y)
	d1y := min(int(math.Ceil(maxY)), cv.Clip.Y+cv.Clip.Height)
	return d0x, d1x, d0y, d1y, d0x < d1x && d0y < d1y
}
