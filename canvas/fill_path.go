package canvas

import "math"

// FillPath fills the inside of a path with a colour, painted through the
// canvas's clip and blend the way every other draw call is. The rule decides
// which insides count: nonzero for shapes that wind around a point, where two
// nested rings drawn the same way are a ring and a hole, and evenodd for shapes
// that must leave a hole wherever they cross themselves.
func (cv *Canvas) FillPath(path *Path, c Color, rule FillRule) {
	if path == nil || c.A() == 0 {
		return
	}
	minX, minY, maxX, maxY, ok := path.Bounds()
	if !ok {
		return
	}
	// The band is the path's box clipped to the clip region: a path reaching out
	// of the canvas is still filled where the canvas is, and one lying wholly
	// outside it is not filled at all.
	x0 := clampInt(int(math.Floor(minX)), cv.Clip.X, cv.Clip.X+cv.Clip.Width)
	y0 := clampInt(int(math.Floor(minY)), cv.Clip.Y, cv.Clip.Y+cv.Clip.Height)
	x1 := clampInt(int(math.Ceil(maxX)), cv.Clip.X, cv.Clip.X+cv.Clip.Width)
	y1 := clampInt(int(math.Ceil(maxY)), cv.Clip.Y, cv.Clip.Y+cv.Clip.Height)
	if x0 >= x1 || y0 >= y1 {
		return
	}
	for _, s := range fillCoverage(path, x0, y0, x1, y1, rule) {
		cv.fillCoverageSpan(s, c)
	}
}
