package svg

// quad reads a quadratic curve: one control point and the point it ends on.
func (st *pathState) quad() bool {
	c, ok1 := st.sc.point()
	e, ok2 := st.sc.point()
	if !ok1 || !ok2 {
		st.warn("the %c command needs four numbers", st.cmd)
		return false
	}
	if st.rel {
		c, e = add(st.cur, c), add(st.cur, e)
	}
	st.p.QuadTo(c.X, c.Y, e.X, e.Y)
	st.cur, st.ctl, st.prevCubic, st.prevQuad = e, c, false, true
	return true
}

// smoothQuad reads a quadratic curve that leaves its control point out: it is
// the mirror of the last one, and the pen itself when the curve before was not
// a quadratic.
func (st *pathState) smoothQuad() bool {
	e, ok := st.sc.point()
	if !ok {
		st.warn("the %c command needs two numbers", st.cmd)
		return false
	}
	if st.rel {
		e = add(st.cur, e)
	}
	c := st.cur
	if st.prevQuad {
		c = mirror(st.ctl, st.cur)
	}
	st.p.QuadTo(c.X, c.Y, e.X, e.Y)
	st.cur, st.ctl, st.prevCubic, st.prevQuad = e, c, false, true
	return true
}
