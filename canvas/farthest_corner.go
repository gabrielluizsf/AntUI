package canvas

// farthestCorner is the corner the given centre is pushed toward.
func farthestCorner(cx, cy float64, w, h int) (float64, float64) {
	ex, ey := 0.0, 0.0
	if cx >= float64(w)/2 {
		ex = float64(w)
	}
	if cy >= float64(h)/2 {
		ey = float64(h)
	}
	return ex, ey
}
