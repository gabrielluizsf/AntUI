package font

import "math"

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
