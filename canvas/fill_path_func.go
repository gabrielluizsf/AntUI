package canvas

import "math"

// FillPathFunc fills the inside of a path the way [Canvas.FillPath] does, with a
// colour that is decided per pixel instead of once for the whole shape. The shade
// is asked for each pixel of each covered run, so a shape painted by a gradient,
// a mesh or anything else that changes across its own area can be filled without
// the canvas having to know what is producing the colour.
//
// The coverage is measured exactly as it is for a flat fill — the same antialiased
// edge, the same fill rule — and comes back as the alpha of the colour the shade
// answered, so a shape whose gradient fades out at its edge gets a soft edge out
// of it for free. The shade is called once per covered pixel and nothing else: the
// pixels the shape does not cover are never asked about, which is what keeps a
// gradient off the rest of the box around it.
//
// The blend and the clip are the canvas's own, so a filled gradient behaves like
// any other fill, including on a layer, where the transparency survives.
func (cv *Canvas) FillPathFunc(path *Path, rule FillRule, shade func(x, y int) Color) {
	if path == nil || shade == nil {
		return
	}
	minX, minY, maxX, maxY, ok := path.Bounds()
	if !ok {
		return
	}
	x0 := clampInt(int(math.Floor(minX)), cv.Clip.X, cv.Clip.X+cv.Clip.Width)
	y0 := clampInt(int(math.Floor(minY)), cv.Clip.Y, cv.Clip.Y+cv.Clip.Height)
	x1 := clampInt(int(math.Ceil(maxX)), cv.Clip.X, cv.Clip.X+cv.Clip.Width)
	y1 := clampInt(int(math.Ceil(maxY)), cv.Clip.Y, cv.Clip.Y+cv.Clip.Height)
	if x0 >= x1 || y0 >= y1 {
		return
	}
	for _, s := range fillCoverage(path, x0, y0, x1, y1, rule) {
		cv.shadeSpan(s, shade)
	}
}

// shadeSpan paints one run one pixel at a time, because a run that shares a
// coverage is not a run that shares a colour once the colour is coming from a
// gradient: the two ends of a run across a gradient are as far apart as the
// gradient makes them. The coverage folds into the alpha the shade answered, so
// a part-covered edge keeps both the colour of the gradient there and the shape
// of its edge.
func (cv *Canvas) shadeSpan(s span, shade func(x, y int) Color) {
	for x := s.x0; x < s.x1; x++ {
		c := shade(x, s.y)
		a := uint8(float64(c.A())*s.a + 0.5)
		if a == 0 {
			continue
		}
		if a != c.A() {
			c = RGBA(c.R(), c.G(), c.B(), a)
		}
		cv.Pixel(x, s.y, c)
	}
}
