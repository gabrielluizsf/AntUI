package canvas

// FillGradient samples a gradient across a rectangle and writes every pixel
// with no blending, which is what a layer rasterizer that has already decided
// what it wants needs: composite the layer to blend the gradient's alphas into
// the pixels behind it. It honours the clip, and it works on any canvas
// format. A linear gradient's first colour lands on the starting side of its
// angle, so linear-gradient(red, blue) with no angle runs top to bottom.
func (cv *Canvas) FillGradient(x, y, w, h int, g Gradient) {
	if g.Kind > GConic || len(g.Stops) == 0 || w <= 0 || h <= 0 {
		return
	}
	stops := normalizeGradientStops(g.Stops)
	for py := cv.Clip.Y; py < cv.Clip.Y+cv.Clip.Height; py++ {
		if py < y || py >= y+h {
			continue
		}
		for px := cv.Clip.X; px < cv.Clip.X+cv.Clip.Width; px++ {
			if px < x || px >= x+w {
				continue
			}
			var t float64
			switch g.Kind {
			case GLinear:
				t = linearGradientT(px-x, py-y, w, h, g.Angle)
			case GRadial:
				t = radialGradientT(float64(px-x), float64(py-y), w, h, g)
			case GConic:
				t = conicGradientT(float64(px-x), float64(py-y), w, h, g)
			}
			cv.Put(px, py, gradientColorAt(stops, t))
		}
	}
}
