package canvas

// cutRow takes the alpha of a run of pixels down to what the mask keeps of it:
// keep[x-x0] is how much of the pixel at x that is, and a nil keep empties the
// whole run. Only the alpha moves, so a pixel that already holds what it would
// hold is left alone rather than written back unchanged.
func (cv *Canvas) cutRow(y, x0, x1 int, keep []float64) {
	for x := x0; x < x1; x++ {
		c := cv.At(x, y)
		a := 0
		if keep != nil {
			a = int(float64(c.A())*keep[x-x0] + 0.5)
		}
		if a == int(c.A()) {
			continue
		}
		cv.Put(x, y, Fade(c, a))
	}
}
