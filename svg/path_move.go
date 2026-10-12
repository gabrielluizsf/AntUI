package svg

import "github.com/gabrielluizsf/antui/canvas"

// move reads a moveto: the two numbers that set where the writing goes on from,
// and then the command becomes a line, because a pair of numbers written after a
// move is a line and not another move.
func (st *pathState) move() bool {
	x, ok1 := st.sc.number()
	y, ok2 := st.sc.number()
	if !ok1 || !ok2 {
		st.warn("the %c command needs two numbers", st.cmd)
		return false
	}
	if st.rel && st.started {
		x, y = st.cur.X+x, st.cur.Y+y
	}
	st.cur = canvas.Point{X: x, Y: y}
	st.start = st.cur
	st.p.MoveTo(x, y)
	st.started = true
	st.ctl, st.prevCubic, st.prevQuad = st.cur, false, false
	if st.rel {
		st.cmd = 'l'
	} else {
		st.cmd = 'L'
	}
	return true
}
