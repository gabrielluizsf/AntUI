package canvas

// nearestCorner is the corner the given centre is pulled toward.
func nearestCorner(cx, cy float64, w, h int) (float64, float64) {
	ex, ey := float64(w), float64(h)
	if cx < float64(w)/2 {
		ex = 0
	}
	if cy < float64(h)/2 {
		ey = 0
	}
	return ex, ey
}
