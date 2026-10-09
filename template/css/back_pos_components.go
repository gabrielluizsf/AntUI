package css

// backComp is one background-position component: a position and the axis its
// keyword names (0 horizontal, 1 vertical, 2 neither).
type backComp struct {
	pos  BackPosition
	axis uint8
}

// readBackComps folds a token list into its components. A keyword swallows
// the offset that immediately follows it, so "right 20px top 30px" is two
// components rather than four; a bare length is a neither-axis component.
func readBackComps(tokens []string, ctx Units) ([]backComp, bool) {
	var comps []backComp
	i := 0
	for i < len(tokens) {
		tk := tokens[i]
		i++
		var c backComp
		switch tk {
		case "left":
			c = backComp{pos: BackPosition{Edge: BackPosStart}, axis: 0}
		case "right":
			c = backComp{pos: BackPosition{Edge: BackPosEnd}, axis: 0}
		case "top":
			c = backComp{pos: BackPosition{Edge: BackPosStart}, axis: 1}
		case "bottom":
			c = backComp{pos: BackPosition{Edge: BackPosEnd}, axis: 1}
		case "center":
			c = backComp{pos: BackPosition{Edge: BackPosMiddle}, axis: 2}
		default:
			l, err := parseLengthAt(tk, ctx)
			if err != nil {
				return nil, false
			}
			c = backComp{pos: BackPosition{Edge: BackPosOffset, Off: l}, axis: 2}
		}
		if c.axis != 2 && i < len(tokens) {
			if l, err := parseLengthAt(tokens[i], ctx); err == nil {
				c.pos.Off = l
				i++
			}
		}
		comps = append(comps, c)
	}
	if len(comps) == 0 || len(comps) > 4 {
		return nil, false
	}
	return comps, true
}
