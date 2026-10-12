package svg

import "strings"

// readColorMatrix reads a `feColorMatrix`: the type it is written with and the
// values that go with it, defaulting to the twenty-number identity matrix.
func (p *filterPrim) readColorMatrix(e *element, warn func(string, ...any)) bool {
	typ := strings.ToLower(strings.TrimSpace(e.attr("type")))
	if typ == "" {
		typ = "matrix"
	}
	raw := e.attr("values")
	switch typ {
	case "matrix":
		if raw == "" {
			p.matrix = identityMatrix()
			return true
		}
		ns, ok := numbers(raw)
		if !ok || len(ns) != 20 {
			warn("the feColorMatrix values %q are not the twenty numbers the matrix needs, so the primitive is left out", raw)
			return false
		}
		copy(p.matrix[:], ns)
	case "saturate":
		v := 1.0
		if raw != "" {
			var ok bool
			if v, ok = filterAmount(raw); !ok {
				warn("the feColorMatrix saturate %q is not a number this package can read, so the primitive is left out", raw)
				return false
			}
		}
		p.matrix = saturateMatrix(v)
	case "huerotate":
		v := 0.0
		if raw != "" {
			var err error
			if v, err = parseAngle(raw); err != nil {
				warn("the feColorMatrix hueRotate %q is not an angle this package can read, so the primitive is left out", raw)
				return false
			}
		}
		p.matrix = hueRotateMatrix(v)
	case "luminancetoalpha":
		p.matrix = luminanceMatrix()
	default:
		warn("the feColorMatrix type %q is not one this package can read, so the primitive is left out", typ)
		return false
	}
	return true
}
