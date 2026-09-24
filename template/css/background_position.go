package css

import "strings"

// parseBackgroundPosition reads a comma-separated background-position list.
// Each value holds one position per axis; a single value is expanded across
// the other axis the way CSS does, centring it. A malformed value drops the
// whole property.
func parseBackgroundPosition(raw string, ctx Units) ([]BackPos, bool) {
	var out []BackPos
	for _, part := range splitFields(raw, ',') {
		p, ok := parseBackPosTokens(splitTokens(strings.TrimSpace(part)), ctx)
		if !ok {
			return nil, false
		}
		out = append(out, p)
	}
	if len(out) == 0 {
		return nil, false
	}
	return out, true
}

// parseBackPosTokens folds a token list into one [BackPos]. A component is a
// keyword with the offset that immediately follows it, or a bare length:
// "right 20px", "top", "center", "10% 20px". Keywords lock their component
// onto their own axis; bare lengths and centre fill the free axes in order,
// and an axis left over after both components settle is centred — so "left"
// means left centre and "20px" means a 20px hand-off on a vertically centred
// picture.
func parseBackPosTokens(tokens []string, ctx Units) (BackPos, bool) {
	type comp struct {
		pos  BackPosition
		axis uint8 // 0 horizontal, 1 vertical, 2 neither
	}
	const (
		axNone = 2
	)
	var comps []comp
	i := 0
	for i < len(tokens) {
		tk := tokens[i]
		i++
		var c comp
		switch tk {
		case "left":
			c = comp{pos: BackPosition{Edge: BackPosStart}, axis: 0}
		case "right":
			c = comp{pos: BackPosition{Edge: BackPosEnd}, axis: 0}
		case "top":
			c = comp{pos: BackPosition{Edge: BackPosStart}, axis: 1}
		case "bottom":
			c = comp{pos: BackPosition{Edge: BackPosEnd}, axis: 1}
		case "center":
			c = comp{pos: BackPosition{Edge: BackPosMiddle}, axis: axNone}
		default:
			l, err := parseLengthAt(tk, ctx)
			if err != nil {
				return BackPos{}, false
			}
			c = comp{pos: BackPosition{Edge: BackPosOffset, Off: l}, axis: axNone}
		}
		// A keyword swallows the offset that follows it, so "right 20px top
		// 30px" is two components rather than four.
		if c.axis != axNone && i < len(tokens) {
			if l, err := parseLengthAt(tokens[i], ctx); err == nil {
				c.pos.Off = l
				i++
			}
		}
		comps = append(comps, c)
	}
	if len(comps) == 0 || len(comps) > 4 {
		return BackPos{}, false
	}

	var out BackPos
	placed := [2]bool{}
	assign := func(idx int, p BackPosition) bool {
		if placed[idx] {
			return false
		}
		out[idx] = p
		placed[idx] = true
		return true
	}
	// First pass: components whose keyword names an axis go there. A keyword
	// whose own axis is taken falls back onto the other, the corner-duplicate
	// corner case; a component that finds both taken is malformed.
	for _, c := range comps {
		if c.axis == axNone {
			continue
		}
		if c.axis == 0 {
			if !assign(0, c.pos) && !assign(1, c.pos) {
				return BackPos{}, false
			}
		} else if !assign(1, c.pos) && !assign(0, c.pos) {
			return BackPos{}, false
		}
	}
	// Second pass: bare offsets and centre fill the free axes in order.
	for _, c := range comps {
		if c.axis != axNone {
			continue
		}
		if !assign(0, c.pos) && !assign(1, c.pos) {
			return BackPos{}, false
		}
	}
	// An axis the values never named is centred.
	if !placed[0] {
		out[0] = BackPosition{Edge: BackPosMiddle, Off: Zero()}
	}
	if !placed[1] {
		out[1] = BackPosition{Edge: BackPosMiddle, Off: Zero()}
	}
	return out, true
}
