package svg

// inherit takes from another gradient everything it wrote itself, which is what
// an href means: the stops, and any geometry left out. An attribute written here
// is kept, so a gradient can take the stops of another and move its own line.
func (out *gradient) inherit(from *gradient) {
	if len(out.stops) == 0 {
		out.stops = from.stops
	}
	for name := range from.set {
		if out.set[name] {
			continue
		}
		switch name {
		case "gradientunits":
			out.userSpace = from.userSpace
		case "spreadmethod":
			out.spread = from.spread
		case "gradienttransform":
			out.transform, out.hasXform = from.transform, from.hasXform
		case "x1":
			out.x1 = from.x1
		case "y1":
			out.y1 = from.y1
		case "x2":
			out.x2 = from.x2
		case "y2":
			out.y2 = from.y2
		case "cx":
			out.cx = from.cx
		case "cy":
			out.cy = from.cy
		case "r":
			out.r = from.r
		case "fx":
			out.fx = from.fx
		case "fy":
			out.fy = from.fy
		}
		if p, ok := from.pct[name]; ok {
			out.pct[name] = p
		}
		out.set[name] = true
	}
}
