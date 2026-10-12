package svg

// command reads one command of the path data and the numbers it takes, moving
// the pen and adding to the path. It answers false where the command is not one
// a path has or its numbers are missing or malformed, which the caller turns
// into no path at all.
func (st *pathState) command() bool {
	st.rel = st.cmd >= 'a' && st.cmd <= 'z'
	switch upperCmd(st.cmd) {
	case 'M':
		return st.move()
	case 'L':
		return st.line()
	case 'H':
		return st.hline()
	case 'V':
		return st.vline()
	case 'C':
		return st.cubic()
	case 'S':
		return st.smoothCubic()
	case 'Q':
		return st.quad()
	case 'T':
		return st.smoothQuad()
	case 'A':
		return st.arc()
	case 'Z':
		return st.close()
	}
	st.warn("%q is not a command a path has", string(st.cmd))
	return false
}
