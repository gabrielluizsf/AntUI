package font

import "math"

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
