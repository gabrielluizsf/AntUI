package font

import "math"

// raster accumulates edge contributions and then reads them off as
// coverage.
type raster struct {
	w, h int
	a    []float64 // one more column than w, so the last delta has somewhere
}

func newRaster(w, h int) *raster {
	if w <= 0 || h <= 0 {
		return &raster{}
	}
	return &raster{w: w, h: h, a: make([]float64, (w+1)*h)}
}

// line adds one edge, in pixels, with y downwards.
//
// This is the whole rasteriser. Everything else is arranging to call it.
func (r *raster) line(x0, y0, x1, y1 float64) {
	if r.w == 0 || y0 == y1 {
		return
	}
	// Going up counts one way and going down the other, which is what makes
	// a hole inside a letter a hole: the outer loop and the inner loop wind
	// in opposite directions and cancel.
	dir := 1.0
	if y0 > y1 {
		dir = -1
		x0, y0, x1, y1 = x1, y1, x0, y0
	}
	dxdy := (x1 - x0) / (y1 - y0)

	yStart := max(int(math.Floor(y0)), 0)
	yEnd := min(int(math.Ceil(y1)), r.h)

	for y := yStart; y < yEnd; y++ {
		// The part of this scanline the edge actually spans.
		top := max(y0, float64(y))
		bottom := min(y1, float64(y+1))
		height := bottom - top
		if height <= 0 {
			continue
		}
		xTop := x0 + dxdy*(top-y0)
		xBottom := x0 + dxdy*(bottom-y0)
		r.span(y, xTop, xBottom, height*dir)
	}
}

// span puts one scanline's worth of an edge into the accumulator: the edge
// crosses this row between xa and xb and covers height of it.
func (r *raster) span(y int, xa, xb, height float64) {
	if xa > xb {
		xa, xb = xb, xa
	}
	row := r.a[y*(r.w+1) : (y+1)*(r.w+1)]
	xa, xb = math.Max(xa, 0), math.Min(xb, float64(r.w))
	if xb <= 0 || xa >= float64(r.w) {
		// Entirely left of the mask still matters: everything to its right
		// is inside, so the coverage is added at column zero.
		if xb <= 0 {
			row[0] += height
		}
		return
	}

	first, last := int(xa), int(xb)
	if first == last {
		// The whole crossing is inside one pixel. It covers the part of that
		// pixel to the right of the edge's middle.
		mid := (xa + xb) / 2
		covered := 1 - (mid - float64(first))
		row[first] += height * covered
		if first+1 <= r.w {
			row[first+1] += height * (1 - covered)
		}
		return
	}

	// Several pixels: the fraction of the crossing that falls in each one is
	// proportional to how much of the x range is in it.
	inv := 1 / (xb - xa)
	for px := first; px <= last && px < r.w; px++ {
		lo, hi := math.Max(xa, float64(px)), math.Min(xb, float64(px+1))
		if hi <= lo {
			continue
		}
		part := (hi - lo) * inv * height
		// Within this pixel the edge sits at the middle of its own span, so
		// that much of the pixel is covered and the rest carries over.
		mid := (lo + hi) / 2
		covered := 1 - (mid - float64(px))
		row[px] += part * covered
		if px+1 <= r.w {
			row[px+1] += part * (1 - covered)
		}
	}
}

// mask reads the accumulator off as coverage, summing along each row.
func (r *raster) mask() []uint8 {
	if r.w == 0 {
		return nil
	}
	out := make([]uint8, r.w*r.h)
	for y := range r.h {
		row := r.a[y*(r.w+1) : (y+1)*(r.w+1)]
		sum := 0.0
		for x := range r.w {
			sum += row[x]
			// The winding can go above one where contours overlap, and it
			// can go a hair below zero from rounding.
			v := math.Abs(sum)
			if v > 1 {
				v = 1
			}
			out[y*r.w+x] = uint8(v*255 + 0.5)
		}
	}
	return out
}

// quad flattens a quadratic curve into lines.
//
// The number of pieces comes from how far the control point is from the
// straight line between the ends: a curve that is nearly straight gets one
// piece and a tight corner gets many. Fixing the count instead would either
// draw far too many at small sizes or show corners at large ones.
func (r *raster) quad(x0, y0, cx, cy, x1, y1 float64) {
	dx := (x0+x1)/2 - cx
	dy := (y0+y1)/2 - cy
	steps := int(math.Ceil(math.Sqrt(math.Hypot(dx, dy) * 2)))
	steps = min(max(steps, 1), 32)

	px, py := x0, y0
	for i := 1; i <= steps; i++ {
		t := float64(i) / float64(steps)
		u := 1 - t
		qx := u*u*x0 + 2*u*t*cx + t*t*x1
		qy := u*u*y0 + 2*u*t*cy + t*t*y1
		r.line(px, py, qx, qy)
		px, py = qx, qy
	}
}
