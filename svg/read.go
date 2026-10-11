package svg

import (
	"strconv"
	"strings"
)

// read takes the attributes written on the gradient. The defaults are already in
// place, so an attribute that is not there is left alone: that is what lets an
// href decide it afterwards.
func (g *gradient) read(attrs map[string]string, warn func(string, ...any)) {
	units := strings.TrimSpace(attrs["gradientunits"])
	if units != "" {
		g.userSpace = strings.EqualFold(units, "userspaceonuse")
		g.set["gradientunits"] = true
	}
	if v := strings.TrimSpace(attrs["spreadmethod"]); v != "" {
		g.set["spreadmethod"] = true
		switch strings.ToLower(v) {
		case "reflect":
			g.spread = spreadReflect
		case "repeat":
			g.spread = spreadRepeat
		default:
			// An unreadable spread method is the default one, which is what the
			// browser draws rather than leaving the ends of the line blank.
			g.spread = spreadPad
		}
	}
	if v := strings.TrimSpace(attrs["gradienttransform"]); v != "" {
		g.set["gradienttransform"] = true
		if m, ok := parseTransform(v); ok {
			g.transform, g.hasXform = m, true
		} else {
			warn("the gradient transform %q is not one this package can read, so the gradient is not moved", v)
		}
	}
	num := func(name string, dst *float64, axis int) {
		raw, ok := attrs[name]
		if !ok {
			return
		}
		g.set[name] = true
		if body, isPct := strings.CutSuffix(strings.TrimSpace(raw), "%"); isPct {
			// A percentage is kept as it was written, because which of the two
			// things it is a percentage of is decided by the units, and the units
			// may still be coming from an href that has not been followed yet.
			// Only that it is a number is worth saying now.
			if _, err := strconv.ParseFloat(strings.TrimSpace(body), 64); err != nil {
				warn("%s %q is not a number, so it stays at %g", name, raw, *dst)
				return
			}
			g.pct[name] = pctRead{raw: strings.TrimSpace(raw), axis: axis}
			return
		}
		v, ok := gradientLength(raw, 0)
		if !ok {
			warn("%s %q is not a length, so it stays at %g", name, raw, *dst)
			return
		}
		*dst = v
	}
	if g.radial {
		num("cx", &g.cx, 0)
		num("cy", &g.cy, 1)
		num("r", &g.r, -1)
		// The focus follows the centre where it does not say otherwise, which is
		// decided once the centre is settled — and it may only be settled by an
		// href, so it waits for [gradient.resolved].
		num("fx", &g.fx, 0)
		num("fy", &g.fy, 1)
		if g.r < 0 {
			warn("the radius %g is negative, so the gradient has none and is painted as one colour", g.r)
			g.r = 0
		}
		return
	}
	num("x1", &g.x1, 0)
	num("y1", &g.y1, 1)
	num("x2", &g.x2, 0)
	num("y2", &g.y2, 1)
}
