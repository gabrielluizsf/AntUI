package svg

import "github.com/gabrielluizsf/antui/canvas"

// arc reads an elliptical arc: the two radii, the rotation, the two flags and
// the point it ends on. The flags are written as single digits that may run
// into the numbers around them, which is why they are read as flags rather than
// as numbers.
func (st *pathState) arc() bool {
	rx, ok1 := st.sc.number()
	ry, ok2 := st.sc.number()
	rot, ok3 := st.sc.number()
	large, ok4 := st.sc.flag()
	sweep, ok5 := st.sc.flag()
	e, ok6 := st.sc.point()
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 {
		st.warn("the %c command needs seven numbers", st.cmd)
		return false
	}
	if st.rel {
		e = add(st.cur, e)
	}
	st.p.ArcTo(canvas.ArcParams{RX: rx, RY: ry, Rot: rot, Large: large, Sweep: sweep, End: canvas.Point{X: e.X, Y: e.Y}})
	st.cur, st.ctl, st.prevCubic, st.prevQuad = e, st.cur, false, false
	return true
}
