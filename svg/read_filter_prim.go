package svg

import (
	"strconv"
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// readFilterPrim reads one primitive of a `<filter>` — the six this package
// can apply, with the inputs it takes and the name it leaves its result under.
// Anything else inside a `<filter>`, and a primitive whose own numbers are not
// the ones it needs, is said and left out: the rest of the filter still runs,
// the same way one unreadable function in the attribute's list does.
func readFilterPrim(e *element, warn func(string, ...any)) (filterPrim, bool) {
	p := filterPrim{in: e.attr("in"), in2: e.attr("in2"), result: e.attr("result")}
	switch e.Name {
	case "feGaussianBlur":
		p.kind = primGaussianBlur
		if raw := e.attr("stdDeviation"); raw != "" {
			ns, ok := numbers(raw)
			if !ok || len(ns) == 0 {
				warn("the feGaussianBlur stdDeviation %q is not a number this package can read, so the primitive is left out", raw)
				return filterPrim{}, false
			}
			p.stdDev = ns[0]
			p.stdDevY = ns[0]
			if p.stdDev < 0 {
				warn("the feGaussianBlur stdDeviation %q is negative, so the primitive is left out", raw)
				return filterPrim{}, false
			}
			if len(ns) > 1 {
				if ns[1] < 0 {
					warn("the feGaussianBlur stdDeviation %q is negative, so the primitive is left out", raw)
					return filterPrim{}, false
				}
				p.stdDevY = ns[1]
			}
		}
	case "feColorMatrix":
		p.kind = primColorMatrix
		typ := strings.ToLower(strings.TrimSpace(e.attr("type")))
		if typ == "" {
			typ = "matrix"
		}
		raw := e.attr("values")
		switch typ {
		case "matrix":
			if raw == "" {
				p.matrix = identityMatrix()
				break
			}
			ns, ok := numbers(raw)
			if !ok || len(ns) != 20 {
				warn("the feColorMatrix values %q are not the twenty numbers the matrix needs, so the primitive is left out", raw)
				return filterPrim{}, false
			}
			copy(p.matrix[:], ns)
		case "saturate":
			v := 1.0
			if raw != "" {
				var ok bool
				if v, ok = filterAmount(raw); !ok {
					warn("the feColorMatrix saturate %q is not a number this package can read, so the primitive is left out", raw)
					return filterPrim{}, false
				}
			}
			p.matrix = saturateMatrix(v)
		case "huerotate":
			v := 0.0
			if raw != "" {
				var err error
				if v, err = parseAngle(raw); err != nil {
					warn("the feColorMatrix hueRotate %q is not an angle this package can read, so the primitive is left out", raw)
					return filterPrim{}, false
				}
			}
			p.matrix = hueRotateMatrix(v)
		case "luminancetoalpha":
			p.matrix = luminanceMatrix()
		default:
			warn("the feColorMatrix type %q is not one this package can read, so the primitive is left out", typ)
			return filterPrim{}, false
		}
	case "feOffset":
		p.kind = primOffset
		for _, at := range []struct {
			attr string
			slot *float64
		}{{"dx", &p.dx}, {"dy", &p.dy}} {
			if raw := e.attr(at.attr); raw != "" {
				v, ok := parseLength(raw)
				if !ok {
					warn("the feOffset %s %q is not a length this package can read, so the primitive is left out", at.attr, raw)
					return filterPrim{}, false
				}
				*at.slot = v
			}
		}
	case "feFlood":
		p.kind = primFlood
		p.flood = canvas.Color(0xFF000000)
		p.floodOpacity = 1
		if raw := e.attr("flood-color"); raw != "" {
			if isCurrentColor(raw) {
				p.floodCurrent = true
			} else if c, ok := parsePaint(raw); ok {
				p.flood = c
			} else {
				warn("the feFlood flood-color %q is not a colour this package can read, so black is flooded instead", raw)
			}
		}
		if raw := e.attr("flood-opacity"); raw != "" {
			v, ok := filterAmount(raw)
			if !ok {
				warn("the feFlood flood-opacity %q is not a number this package can read, so it is left out", raw)
			} else {
				p.floodOpacity = v
			}
		}
	case "feComposite":
		p.kind = primComposite
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
					return filterPrim{}, false
				}
				p.k[i] = v
			}
		}
	case "feBlend":
		p.kind = primBlend
		mode := strings.ToLower(strings.TrimSpace(e.attr("mode")))
		if mode == "" {
			mode = "normal"
		}
		switch mode {
		case "normal":
			p.blend = blendNormal
		case "multiply":
			p.blend = blendMultiply
		case "screen":
			p.blend = blendScreen
		case "darken":
			p.blend = blendDarken
		case "lighten":
			p.blend = blendLighten
		default:
			warn("the feBlend mode %q is not one this package can read, so \"normal\" is used instead", mode)
		}
	default:
		warn("the filter primitive %q is not one this package can read, so it is left out", e.Name)
		return filterPrim{}, false
	}
	return p, true
}
