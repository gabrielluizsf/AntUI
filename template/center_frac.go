package template

import "github.com/gabrielluizsf/antui/template/css"

// centerFrac is where the centre of a gradient lands as a fraction of one
// axis. A bare offset runs from the starting edge; a right or bottom anchor
// ranges back towards the inside.
func (cs *cssStyle) centerFrac(st css.Style, pos css.BackPosition, u, dim int) float64 {
	frac := 0.5
	switch pos.Edge {
	case css.BackPosEnd:
		frac = 1
	case css.BackPosStart, css.BackPosOffset:
		frac = 0
	}
	den := float64(max(dim, 1))
	var off float64
	if pos.Off.IsPct() {
		off = float64(pos.Off.Resolve(css.Units{Width: dim})) / den
	} else {
		off = float64(pos.Off.Resolve(css.Units{Width: dim})*u) / den
	}
	if pos.Edge == css.BackPosEnd {
		return frac - off
	}
	return frac + off
}
