package canvas

// maskApply paints one row of a shape mask onto a canvas: the union of what
// every measured shape covers, with everything outside the band emptied.
type maskApply struct {
	cv       *Canvas
	all      []maskMeasured
	runs     *maskRuns
	row      []float64
	cx0, cx1 int
	bx0, bx1 int
	by0, by1 int
}

func (a *maskApply) paint(y int) {
	if y < a.by0 || y >= a.by1 {
		a.cv.cutRow(y, a.cx0, a.cx1, nil)
		return
	}
	clear(a.row)
	for i := range a.all {
		m := a.runs.cov(&a.all[i], 0, y)
		for x := range a.row {
			// The most of a pixel any one shape covers is how much of
			// it the union covers: a second shape over the same pixel
			// adds nothing to what is already inside the clip.
			if m[x] > a.row[x] {
				a.row[x] = m[x]
			}
		}
	}
	// The two ends of the row are past the band and so past every shape;
	// the middle of it keeps what the shapes measured.
	a.cv.cutRow(y, a.cx0, a.bx0, nil)
	a.cv.cutRow(y, a.bx0, a.bx1, a.row)
	a.cv.cutRow(y, a.bx1, a.cx1, nil)
}
