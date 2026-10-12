package svg

// close reads a closepath: the subpath is joined back to where it began, and the
// pen stands there for the next command, with no command of its own left over.
func (st *pathState) close() bool {
	st.p.Close()
	st.cur = st.start
	st.started = true
	st.ctl, st.prevCubic, st.prevQuad = st.cur, false, false
	st.cmd = 0
	return true
}
