package css

// bezierX and bezierY evaluate the cubic applied to a parameter u in [0,1].
func bezierX(t Timing, u float64) float64 {
	return 3*(1-u)*(1-u)*u*t.x1 + 3*(1-u)*u*u*t.x2 + u*u*u
}

func bezierY(t Timing, u float64) float64 {
	return 3*(1-u)*(1-u)*u*t.y1 + 3*(1-u)*u*u*t.y2 + u*u*u
}
