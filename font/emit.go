package font

// Emit walks one contour, turning the on/off-curve points into lines and
// quadratics.
//
// TrueType leaves the on-curve points between two off-curve ones out, and
// they are the midpoint — so a run of control points is a run of curves that
// meet smoothly, written down as half the data.
func Emit(r *raster, c Contour, scale, ox, oy float64) {
	n := len(c.X)
	if n < 2 {
		return
	}

	at := func(i int) (float64, float64, bool) {
		i = ((i % n) + n) % n
		return c.X[i]*scale - ox, oy - c.Y[i]*scale, c.On[i]
	}

	start := -1
	for i := range n {
		if c.On[i] {
			start = i
			break
		}
	}
	var sx, sy float64
	if start < 0 {
		ax, ay, _ := at(0)
		bx, by, _ := at(1)
		sx, sy = (ax+bx)/2, (ay+by)/2
		start = -1
	} else {
		sx, sy, _ = at(start)
	}

	x, y := sx, sy
	for k := 1; k <= n; k++ {
		cx, cy, on := at(start + k)
		if on {
			r.line(x, y, cx, cy)
			x, y = cx, cy
			continue
		}

		nx, ny, nOn := at(start + k + 1)
		if nOn {
			k++ 
		} else {
			nx, ny = (cx+nx)/2, (cy+ny)/2
		}
		r.quad(x, y, cx, cy, nx, ny)
		x, y = nx, ny
	}
	r.line(x, y, sx, sy) 
}
