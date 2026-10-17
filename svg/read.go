package svg

import "strings"

// read takes the attributes written on the gradient. The defaults are already in
// place, so an attribute that is not there is left alone: that is what lets an
// href decide it afterwards.
func (g *gradient) read(attrs map[string]string, warn func(string, ...any)) {
	units := strings.TrimSpace(attrs["gradientUnits"])
	if units != "" {
		g.userSpace = strings.EqualFold(units, "userspaceonuse")
		g.set["gradientUnits"] = true
	}
	if v := strings.TrimSpace(attrs["spreadMethod"]); v != "" {
		g.set["spreadMethod"] = true
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
	if v := strings.TrimSpace(attrs["gradientTransform"]); v != "" {
		g.set["gradientTransform"] = true
		if m, ok := parseTransform(v); ok {
			g.transform, g.hasXform = m, true
		} else {
			warn("the gradient transform %q is not one this package can read, so the gradient is not moved", v)
		}
	}
	if g.radial {
		g.readNum(attrs, "cx", &g.cx, 0, warn)
		g.readNum(attrs, "cy", &g.cy, 1, warn)
		g.readNum(attrs, "r", &g.r, -1, warn)
		// The focus follows the centre where it does not say otherwise, which is
		// decided once the centre is settled — and it may only be settled by an
		// href, so it waits for [gradient.resolved].
		g.readNum(attrs, "fx", &g.fx, 0, warn)
		g.readNum(attrs, "fy", &g.fy, 1, warn)
		if g.r < 0 {
			warn("the radius %g is negative, so the gradient has none and is painted as one colour", g.r)
			g.r = 0
		}
		return
	}
	g.readNum(attrs, "x1", &g.x1, 0, warn)
	g.readNum(attrs, "y1", &g.y1, 1, warn)
	g.readNum(attrs, "x2", &g.x2, 0, warn)
	g.readNum(attrs, "y2", &g.y2, 1, warn)
}
