package svg

import (
	"strconv"
	"strings"
)

// readComposite reads a `feComposite`: the operator, which is `over` where
// none is written, and the four k values arithmetic needs.
func (p *filterPrim) readComposite(e *element, warn func(string, ...any)) bool {
	op := strings.ToLower(strings.TrimSpace(e.attr("operator")))
	if op == "" {
		op = "over"
	}
	switch op {
	case "over":
		p.comp = compOver
	case "in":
		p.comp = compIn
	case "out":
		p.comp = compOut
	case "atop":
		p.comp = compAtop
	case "xor":
		p.comp = compXor
	case "lighter":
		p.comp = compLighter
	case "arithmetic":
		p.comp = compArithmetic
	default:
		warn("the feComposite operator %q is not one this package can read, so \"over\" is used instead", op)
	}
	for i, at := range []string{"k1", "k2", "k3", "k4"} {
		if raw := e.attr(at); raw != "" {
			v, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
			if err != nil {
				warn("the feComposite %s %q is not a number this package can read, so the primitive is left out", at, raw)
				return false
			}
			p.k[i] = v
		}
	}
	return true
}
