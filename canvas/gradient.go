package canvas

import "math"

// GradientKind names one CSS gradient the canvas can fill a rectangle with.
// The values are deliberately independent of any stylesheet; a template maps
// its own parsed gradients onto them.
type GradientKind uint8

// The gradient families.
const (
	GLinear GradientKind = iota
	GRadial
	GConic
)

// GradientShape is the radial ending shape keyword.
type GradientShape uint8

// The radial shapes.
const (
	GradEllipse GradientShape = iota
	GradCircle
)

// GradientSize is the radial size keyword that fixes how far the ending shape
// reaches.
type GradientSize uint8

// The radial sizes, from the CSS farthest-corner default.
const (
	GSFarthestCorner GradientSize = iota
	GSClosestSide
	GSFarthestSide
	GSClosestCorner
)

// GradientStop is one colour stop of a gradient: a colour and where along the
// gradient line it changes, running from 0 to 1. The offsets a caller leaves
// out are spread evenly by [Canvas.FillGradient].
type GradientStop struct {
	Offset float64
	Color  Color
}

// Gradient is the image a template paints into a rectangle with
// [Canvas.FillGradient]. Angle is radians measured clockwise from the top,
// matching "to top" being 0°; CenterX and CenterY are fractions of the
// rectangle the radial or conic gradient spins around. An Angle of zero and a
// centre of 0.5, 0.5 are the defaults a gradient without them wants.
type Gradient struct {
	Kind    GradientKind
	Angle   float64
	CenterX float64
	CenterY float64
	Shape   GradientShape
	Size    GradientSize
	Stops   []GradientStop
}

// FillGradient samples a gradient across a rectangle and writes every pixel
// with no blending, which is what a layer rasterizer that has already decided
// what it wants needs: composite the layer to blend the gradient's alphas into
// the pixels behind it. It honours the clip, and it works on any canvas
// format. A linear gradient's first colour lands on the starting side of its
// angle, so linear-gradient(red, blue) with no angle runs top to bottom.
func (cv *Canvas) FillGradient(x, y, w, h int, g Gradient) {
	if g.Kind > GConic || len(g.Stops) == 0 || w <= 0 || h <= 0 {
		return
	}
	stops := normalizeGradientStops(g.Stops)
	for py := cv.Clip.Y; py < cv.Clip.Y+cv.Clip.Height; py++ {
		if py < y || py >= y+h {
			continue
		}
		for px := cv.Clip.X; px < cv.Clip.X+cv.Clip.Width; px++ {
			if px < x || px >= x+w {
				continue
			}
			var t float64
			switch g.Kind {
			case GLinear:
				t = linearGradientT(px-x, py-y, w, h, g.Angle)
			case GRadial:
				t = radialGradientT(float64(px-x), float64(py-y), w, h, g)
			case GConic:
				t = conicGradientT(float64(px-x), float64(py-y), w, h, g)
			}
			cv.Put(px, py, gradientColorAt(stops, t))
		}
	}
}

// normalizeGradientStops fills the offsets a caller left out the way CSS does:
// the first and last missing stops become 0 and 1, and missing stops between
// known ones spread evenly over the gap.
func normalizeGradientStops(in []GradientStop) []GradientStop {
	n := len(in)
	if n == 0 {
		return nil
	}
	if n == 1 {
		return []GradientStop{{Offset: 0, Color: in[0].Color}, {Offset: 1, Color: in[0].Color}}
	}
	out := make([]GradientStop, n)
	copy(out, in)
	has := make([]bool, n)
	any := false
	for i := range n {
		has[i] = out[i].Offset >= 0
		any = any || has[i]
	}
	if !any {
		for i := range n {
			out[i].Offset = float64(i) / float64(n-1)
		}
		return out
	}
	if !has[0] {
		out[0].Offset = 0
		has[0] = true
	}
	if !has[n-1] {
		out[n-1].Offset = 1
		has[n-1] = true
	}
	// Fill each run of unset stops between known neighbours.
	i := 0
	for i < n {
		if has[i] {
			i++
			continue
		}
		j := i
		for j < n && !has[j] {
			j++
		}
		lo, hi := out[i-1].Offset, out[j].Offset
		for k := i; k < j; k++ {
			f := float64(k-i+1) / float64(j-i+1)
			out[k].Offset = lo + (hi-lo)*f
		}
		i = j
	}
	return out
}

// gradientColorAt interpolates the colour at one point along a gradient line.
// Points outside the stops clamp to the end colours.
func gradientColorAt(stops []GradientStop, t float64) Color {
	if t <= 0 {
		return stops[0].Color
	}
	if t >= 1 {
		return stops[len(stops)-1].Color
	}
	for i := 1; i < len(stops); i++ {
		if t < stops[i].Offset {
			a, b := stops[i-1], stops[i]
			span := b.Offset - a.Offset
			var f float32
			if span > 0 {
				f = float32((t - a.Offset) / span)
			}
			return Mix(a.Color, b.Color, f)
		}
	}
	return stops[len(stops)-1].Color
}

// linearGradientT is where a point lands along the gradient line: 0 at the
// start of the angle's direction, 1 at the far edge. The line runs through
// the rectangle's centre and its extent is set by projecting the four corners
// onto it, which keeps any angle inside the box and puts the corner
// directions exactly corner to corner.
func linearGradientT(x, y, w, h int, angle float64) float64 {
	dx := math.Sin(angle)
	dy := -math.Cos(angle)
	cw, ch := float64(w)/2, float64(h)/2
	minP, maxP := math.Inf(1), math.Inf(-1)
	for _, c := range [4][2]float64{{0, 0}, {float64(w), 0}, {0, float64(h)}, {float64(w), float64(h)}} {
		p := (c[0]-cw)*dx + (c[1]-ch)*dy
		if p < minP {
			minP = p
		}
		if p > maxP {
			maxP = p
		}
	}
	if span := maxP - minP; span > 0 {
		return ((float64(x)-cw)*dx + (float64(y)-ch)*dy - minP) / span
	}
	return 0
}

// radialGradientT scales the distance from the centre up to the ending shape
// into a point along the gradient line. The ending shape's size keywords are
// resolved against the rectangle, and a circle uses the largest circle that
// the keywords name while an ellipse keeps its two radii live.
func radialGradientT(x, y float64, w, h int, g Gradient) float64 {
	cx := g.CenterX * float64(w)
	cy := g.CenterY * float64(h)
	if g.Shape == GradCircle {
		r := circleRadius(cx, cy, w, h, g.Size)
		if r <= 0 {
			return 0
		}
		px, py := x-cx, y-cy
		return math.Sqrt(px*px+py*py) / r
	}
	var rx, ry float64
	switch g.Size {
	case GSClosestSide:
		rx = min(cx, float64(w)-cx)
		ry = min(cy, float64(h)-cy)
	case GSFarthestSide:
		rx = max(cx, float64(w)-cx)
		ry = max(cy, float64(h)-cy)
	case GSClosestCorner:
		ex, ey := nearestCorner(cx, cy, w, h)
		rx, ry = math.Abs(ex-cx), math.Abs(ey-cy)
	default:
		ex, ey := farthestCorner(cx, cy, w, h)
		rx, ry = math.Abs(ex-cx), math.Abs(ey-cy)
	}
	if rx <= 0 || ry <= 0 {
		return 0
	}
	px, py := x-cx, y-cy
	return math.Sqrt(px*px/(rx*rx) + py*py/(ry*ry))
}

// circleRadius is how far a circular ending shape of the given size reaches.
func circleRadius(cx, cy float64, w, h int, size GradientSize) float64 {
	switch size {
	case GSClosestSide:
		return min(min(cx, float64(w)-cx), min(cy, float64(h)-cy))
	case GSFarthestSide:
		return max(max(cx, float64(w)-cx), max(cy, float64(h)-cy))
	case GSClosestCorner:
		ex, ey := nearestCorner(cx, cy, w, h)
		return math.Hypot(ex-cx, ey-cy)
	default:
		ex, ey := farthestCorner(cx, cy, w, h)
		return math.Hypot(ex-cx, ey-cy)
	}
}

// farthestCorner is the corner the given centre is pushed toward.
func farthestCorner(cx, cy float64, w, h int) (float64, float64) {
	ex, ey := 0.0, 0.0
	if cx >= float64(w)/2 {
		ex = float64(w)
	}
	if cy >= float64(h)/2 {
		ey = float64(h)
	}
	return ex, ey
}

// nearestCorner is the corner the given centre is pulled toward.
func nearestCorner(cx, cy float64, w, h int) (float64, float64) {
	ex, ey := float64(w), float64(h)
	if cx < float64(w)/2 {
		ex = 0
	}
	if cy < float64(h)/2 {
		ey = 0
	}
	return ex, ey
}

// conicGradientT walks one revolution around the centre, 0 at the start angle
// (which points up by default) and 1 back around clockwise.
func conicGradientT(x, y float64, w, h int, g Gradient) float64 {
	cx := g.CenterX * float64(w)
	cy := g.CenterY * float64(h)
	ang := math.Atan2(x-cx, -(y-cy)) - g.Angle
	ang = math.Mod(ang, 2*math.Pi)
	if ang < 0 {
		ang += 2 * math.Pi
	}
	return ang / (2 * math.Pi)
}
