package svg

import "math"

// resolvePercentages is where the percentages kept as written become numbers,
// now that the units they are a percentage of are finally known. A radius lands
// on the normalized diagonal of the drawing in user units, which is what SVG
// measures it against where it has no side of its own to take.
func (out *gradient) resolvePercentages() {
	for name, p := range out.pct {
		v, ok := gradientLength(p.raw, out.percentOf(p.axis))
		if !ok {
			continue
		}
		switch name {
		case "x1":
			out.x1 = v
		case "y1":
			out.y1 = v
		case "x2":
			out.x2 = v
		case "y2":
			out.y2 = v
		case "cx":
			out.cx = v
		case "cy":
			out.cy = v
		case "r":
			out.r = clampFloat(v, 0, math.Inf(1))
		case "fx":
			out.fx = v
		case "fy":
			out.fy = v
		}
	}
}
