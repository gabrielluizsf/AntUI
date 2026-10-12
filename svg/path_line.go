package svg

import "github.com/gabrielluizsf/antui/canvas"

// line reads a lineto, straight to the point the two numbers give.
func (st *pathState) line() bool {
	x, ok1 := st.sc.number()
	y, ok2 := st.sc.number()
	if !ok1 || !ok2 {
		st.warn("the %c command needs two numbers", st.cmd)
		return false
	}
	if st.rel {
		x, y = st.cur.X+x, st.cur.Y+y
	}
	st.cur = canvas.Point{X: x, Y: y}
	st.p.LineTo(x, y)
	st.ctl, st.prevCubic, st.prevQuad = st.cur, false, false
	return true
}

// hline reads a horizontal lineto, one number for where across it goes.
func (st *pathState) hline() bool {
	x, ok := st.sc.number()
	if !ok {
		st.warn("the %c command needs one number", st.cmd)
		return false
	}
	if st.rel {
		x = st.cur.X + x
	}
	st.cur = canvas.Point{X: x, Y: st.cur.Y}
	st.p.LineTo(x, st.cur.Y)
	st.ctl, st.prevCubic, st.prevQuad = st.cur, false, false
	return true
}

// vline reads a vertical lineto, one number for how far down it goes.
func (st *pathState) vline() bool {
	y, ok := st.sc.number()
	if !ok {
		st.warn("the %c command needs one number", st.cmd)
		return false
	}
	if st.rel {
		y = st.cur.Y + y
	}
	st.cur = canvas.Point{X: st.cur.X, Y: y}
	st.p.LineTo(st.cur.X, y)
	st.ctl, st.prevCubic, st.prevQuad = st.cur, false, false
	return true
}
