package svg

import "math"

// circleFrom is how far along a ray from inside a circle that ray travels before
// it leaves the circle, which is the distance from the focus to the far side of
// it as seen from where the ray is pointing.
func circleFrom(o point, ux, uy float64, c point, r float64) (float64, bool) {
	fx, fy := o.X-c.X, o.Y-c.Y
	b := 2 * (ux*fx + uy*fy)
	disc := b*b - 4*(fx*fx+fy*fy-r*r)
	if disc <= 0 {
		return 0, false
	}
	// The far root is the one the ray reaches on its way out rather than the one
	// behind it, and the ray starts at the focus, so its distance along is the
	// positive one.
	t := (-b + math.Sqrt(disc)) / 2
	if t <= 0 {
		return 0, false
	}
	return t, true
}
