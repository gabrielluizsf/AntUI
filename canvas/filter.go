package canvas

// FilterKind names one CSS filter function the canvas can apply to a region.
// The values are deliberately independent of any stylesheet; a template maps
// its own parsed filters onto them.
type FilterKind uint8

// The filter functions. Blur takes a radius in pixels; hue-rotate takes
// degrees; the rest take a fraction, where 1 is the identity for contrast and
// brightness and the full effect for grayscale, sepia and invert.
const (
	FilterGrayscale FilterKind = iota
	FilterSepia
	FilterInvert
	FilterBrightness
	FilterContrast
	FilterHueRotate
	FilterBlur
)

// FilterRegion applies one filter to the rectangle in place. It works on an
// ARGB32 canvas; the narrow formats have no per-pixel colour to filter and are
// left untouched. A blur is handled by [Canvas.Blur] and is the one filter
// that reads outside the rectangle, so callers give it room to bleed.
func (cv *Canvas) FilterRegion(x, y, w, h int, kind FilterKind, amount float64) {
	if cv.Pixels == nil || w <= 0 || h <= 0 {
		return
	}
	if kind == FilterBlur {
		cv.Blur(x, y, w, h, int(amount+0.5))
		return
	}
	x0, y0, x1, y1 := cv.filterBounds(x, y, w, h)
	cv.markDirtyRect(Area{X: x0, Y: y0, Width: x1 - x0, Height: y1 - y0})
	for py := y0; py < y1; py++ {
		row := cv.Pixels[py*cv.Stride : py*cv.Stride+cv.Width]
		for px := x0; px < x1; px++ {
			row[px] = applyFilter(row[px], kind, amount)
		}
	}
}

// filterBounds intersects a rectangle with the canvas and its active clip.
func (cv *Canvas) filterBounds(x, y, w, h int) (x0, y0, x1, y1 int) {
	x0 = max(x, cv.Clip.X)
	y0 = max(y, cv.Clip.Y)
	x1 = min(x+w, cv.Clip.X+cv.Clip.Width)
	y1 = min(y+h, cv.Clip.Y+cv.Clip.Height)
	x0 = max(x0, 0)
	y0 = max(y0, 0)
	x1 = min(x1, cv.Width)
	y1 = min(y1, cv.Height)
	if x0 >= x1 || y0 >= y1 {
		return 0, 0, 0, 0
	}
	return x0, y0, x1, y1
}

// applyFilter runs one filter on one colour, leaving its alpha alone.
func applyFilter(c Color, kind FilterKind, amount float64) Color {
	switch kind {
	case FilterGrayscale:
		return lerpMat(identityMat, grayscaleMat, clamp01(amount)).apply(c)
	case FilterSepia:
		return lerpMat(identityMat, sepiaMat, clamp01(amount)).apply(c)
	case FilterInvert:
		a := clamp01(amount)
		f := func(v uint8) uint8 { return clampByte(float64(v) + a*(255-2*float64(v))) }
		return RGBA(f(c.R()), f(c.G()), f(c.B()), c.A())
	case FilterBrightness:
		a := max(amount, 0)
		f := func(v uint8) uint8 { return clampByte(float64(v) * a) }
		return RGBA(f(c.R()), f(c.G()), f(c.B()), c.A())
	case FilterContrast:
		a := max(amount, 0)
		f := func(v uint8) uint8 { return clampByte((float64(v)-127.5)*a + 127.5) }
		return RGBA(f(c.R()), f(c.G()), f(c.B()), c.A())
	case FilterHueRotate:
		return hueRotateMat(amount).apply(c)
	}
	return c
}

// clamp01 keeps a filter fraction inside [0,1].
func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

// clampByte rounds a channel and keeps it inside 0..255.
func clampByte(v float64) uint8 {
	if v <= 0 {
		return 0
	}
	if v >= 255 {
		return 255
	}
	return uint8(v + 0.5)
}
