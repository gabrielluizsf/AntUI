package svg

import "github.com/gabrielluizsf/antui/canvas"

// cubic reads a cubic curve: two control points and the point it ends on.
func (st *pathState) cubic() bool {
	c1, ok1 := st.sc.point()
	c2, ok2 := st.sc.point()
	e, ok3 := st.sc.point()
	if !ok1 || !ok2 || !ok3 {
		st.warn("the %c command needs six numbers", st.cmd)
		return false
	}
	if st.rel {
		c1, c2, e = add(st.cur, c1), add(st.cur, c2), add(st.cur, e)
	}
	st.p.CubicTo(canvas.Point{X:c1.X,Y:c1.Y}, canvas.Point{X:c2.X,Y:c2.Y}, canvas.Point{X:e.X,Y:e.Y})
	st.cur, st.ctl, st.prevCubic, st.prevQuad = e, c2, true, false
	return true
}

// smoothCubic reads a cubic curve that leaves its first control point out: it
// is the mirror of the last one, and the pen itself when the curve before was
// not a cubic.
func (st *pathState) smoothCubic() bool {
	c2, ok1 := st.sc.point()
	e, ok2 := st.sc.point()
	if !ok1 || !ok2 {
		st.warn("the %c command needs four numbers", st.cmd)
		return false
	}
	if st.rel {
		c2, e = add(st.cur, c2), add(st.cur, e)
	}
	c1 := st.cur
	if st.prevCubic {
		c1 = mirror(st.ctl, st.cur)
	}
	st.p.CubicTo(canvas.Point{X:c1.X,Y:c1.Y}, canvas.Point{X:c2.X,Y:c2.Y}, canvas.Point{X:e.X,Y:e.Y})
	st.cur, st.ctl, st.prevCubic, st.prevQuad = e, c2, true, false
	return true
}
