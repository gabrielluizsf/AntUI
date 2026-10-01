package svg

import (
	"math"

	"github.com/gabrielluizsf/antui/canvas"
)

// sampler is a gradient with everything that does not depend on the pixel worked
// out once: where its line or its circle is, the stops with the offsets they left
// out filled in, and its spread. A pixel then costs a matrix, a distance and a
// mix, which is what makes a gradient in a drawing that is painted every frame
// affordable at all.
type sampler struct {
	// toGradient is what turns a point of the drawing's own coordinates into a
	// point of the gradient's: gradientTransform undone, and the shape's box
	// opened out to the unit square when the geometry is written in fractions of
	// it. That is the order SVG asks for, and reversing it is what makes a
	// gradientTransform mean the same thing in both of the units.
	toGradient canvas.Matrix
	// stops are the colour changes with the offsets they left out filled in.
	stops []gradientStop
	// spread is what happens past either end of the line.
	spread spread
	// linear says a line running from p0 along d; the others are a circle.
	linear bool
	p0, d  point
	// The circle, and the focus it is brightest at, which is never outside it.
	focus point
	c     point
	r     float64
}

// newSampler is a gradient ready to be asked about the points of the drawing,
// for a shape whose own box is the given one. A shape with no area, a gradient
// with no stops, a gradientTransform that cannot be undone and a radial gradient
// with no radius all answer a sampler with no stops, which paints nothing: those
// are the four ways SVG says a gradient does not paint, rather than four ways of
// dividing by zero.
func newSampler(g *gradient, minX, minY, w, h float64) sampler {
	s := sampler{
		toGradient: canvas.Identity(),
		stops:      normalizeStops(g.stops),
		spread:     g.spread,
		linear:     !g.radial,
		p0:         point{X: g.x1, Y: g.y1},
		d:          point{X: g.x2 - g.x1, Y: g.y2 - g.y1},
		focus:      point{X: g.fx, Y: g.fy},
		c:          point{X: g.cx, Y: g.cy},
		r:          g.r,
	}
	if len(s.stops) == 0 {
		return sampler{}
	}
	if g.hasXform {
		inv, ok := g.transform.Inverse()
		if !ok {
			return sampler{}
		}
		s.toGradient = inv
	}
	if !g.userSpace {
		// The geometry is written in fractions of the shape's own box, so the
		// point is measured from the corner of that box and divided by its size
		// before the gradientTransform is undone in the unit square: that
		// transform belongs to the gradient and is written against the square the
		// fractions are of, not against wherever the shape happens to be. The
		// point goes through the matrix on the right first, so the move to the
		// corner comes before the division and the two of them before the
		// gradientTransform is undone.
		if w <= 0 || h <= 0 {
			return sampler{}
		}
		box := canvas.Scale(1/w, 1/h).Mul(canvas.Translate(-minX, -minY))
		s.toGradient = s.toGradient.Mul(box)
	}
	if g.radial {
		s.focus = insideCircle(s.focus, s.c, s.r)
	}
	return s
}

// colorAt is the colour of the gradient at one point of the drawing, in the same
// coordinates the shape itself was written in.
func (s sampler) colorAt(x, y float64, current canvas.Color) canvas.Color {
	if len(s.stops) == 0 {
		return 0
	}
	gx, gy := s.toGradient.Map(x, y)
	return s.stopColor(s.spreadAt(s.parameter(gx, gy)), current)
}

// parameter is where a point of the gradient's own space lands along its line,
// 0 at the start and 1 at the far end. Past either end it keeps going, because
// what happens out there is the spread method's business and not this one's.
func (s sampler) parameter(x, y float64) float64 {
	if !s.linear {
		return s.radialParameter(x, y)
	}
	span := s.d.X*s.d.X + s.d.Y*s.d.Y
	if span == 0 {
		// A gradient whose ends are in the same place paints its last stop
		// everywhere, which is what SVG says rather than leaving it blank.
		return 1
	}
	return ((x-s.p0.X)*s.d.X + (y-s.p0.Y)*s.d.Y) / span
}

// radialParameter is where a point lands between the focus of a radial gradient,
// where its first stop is, and the circle the last stop is on. The two are not
// the same distance apart from a point that is off to one side, so the point is
// measured against where the ray from the focus through it leaves the circle
// rather than against the radius — which is what gives a focus away from the
// centre its highlight instead of a ring around the middle.
func (s sampler) radialParameter(x, y float64) float64 {
	if s.r <= 0 {
		// A gradient with no radius is a single colour, which is the last one.
		return 1
	}
	dx, dy := x-s.focus.X, y-s.focus.Y
	d := math.Hypot(dx, dy)
	if d == 0 {
		return 0
	}
	if span, ok := circleFrom(s.focus, dx/d, dy/d, s.c, s.r); ok {
		return d / span
	}
	// The ray misses the circle, which cannot happen with a focus inside it but
	// can with a radius so small the arithmetic does not hold: the distance from
	// the focus alone is as close as the point can be placed.
	return d / s.r
}

// circleFrom is how far along a ray from inside a circle that ray travels before
// it leaves the circle, which is the distance from the focus to the far side of
// it as seen from where the ray is pointing.
func circleFrom(o point, ux, uy float64, c point, r float64) (float64, bool) {
	fx, fy := o.X-c.X, o.Y-c.Y
	b := 2 * (ux*fx + uy*fy)
	disc := b*b - 4*(fx*fx+fy*fy-r*r)
	if disc <= 0 {
		return 0, false
	}
	// The far root is the one the ray reaches on its way out rather than the one
	// behind it, and the ray starts at the focus, so its distance along is the
	// positive one.
	t := (-b + math.Sqrt(disc)) / 2
	if t <= 0 {
		return 0, false
	}
	return t, true
}

// insideCircle pulls a focus back inside the circle it belongs to. A focus outside
// it has no highlight to give, and SVG puts it on the edge rather than refusing
// to paint the shape at all.
func insideCircle(f, c point, r float64) point {
	d := math.Hypot(f.X-c.X, f.Y-c.Y)
	if d <= r || d == 0 {
		return f
	}
	return point{X: c.X + (f.X-c.X)*r/d, Y: c.Y + (f.Y-c.Y)*r/d}
}

// spreadAt is the parameter after the spread method has had its say about
// everything past the ends of the line.
func (s sampler) spreadAt(t float64) float64 {
	switch s.spread {
	case spreadRepeat:
		// Mod keeps the ends out of it, which is what makes a repeated gradient
		// start again at the first colour instead of jumping to the last one.
		return t - math.Floor(t)
	case spreadReflect:
		t = math.Mod(math.Abs(t), 2)
		if t > 1 {
			return 2 - t
		}
		return t
	}
	return t
}

// stopColor is the colour at one point along the line, with the ends clamping to
// the first and last stop — which is what the pad method is, and what the other
// two land inside by the time they get here.
func (s sampler) stopColor(t float64, current canvas.Color) canvas.Color {
	stops := s.stops
	switch {
	case t <= 0:
		return stops[0].at(current)
	case t >= 1:
		return stops[len(stops)-1].at(current)
	}
	for i := 1; i < len(stops); i++ {
		if t >= stops[i].offset {
			continue
		}
		a, b := stops[i-1], stops[i]
		var f float32
		if span := b.offset - a.offset; span > 0 {
			f = float32((t - a.offset) / span)
		}
		return canvas.Mix(a.at(current), b.at(current), f)
	}
	return stops[len(stops)-1].at(current)
}

// at is the colour the stop is painted in, which for a stop written as
// `currentColor` is the colour of the painting rather than the one the file named.
func (s gradientStop) at(current canvas.Color) canvas.Color {
	if s.current {
		return current
	}
	return s.color
}

// normalizeStops is the same rule canvas.FillGradient uses for the offsets a
// gradient leaves out: the first and last missing ones become 0 and 1, and the
// missing ones in between spread evenly over the gap they sit in. On top of
// that the offsets are held inside the line and never go backwards, since a stop
// behind the one before it is not a gradient but a piece of the file out of
// order.
func normalizeStops(in []gradientStop) []gradientStop {
	n := len(in)
	if n == 0 {
		return nil
	}
	out := make([]gradientStop, n)
	copy(out, in)
	if n == 1 {
		// One stop is a flat colour, which is what a gradient of one colour is.
		out = append(out, gradientStop{offset: 1, color: in[0].color, current: in[0].current})
		return out
	}
	if !out[0].hasOffset {
		out[0].offset = 0
	}
	if !out[n-1].hasOffset {
		out[n-1].offset = 1
	}
	for i := 0; i < n; {
		if out[i].hasOffset {
			i++
			continue
		}
		j := i
		for j < n && !out[j].hasOffset {
			j++
		}
		lo, hi := out[i-1].offset, out[j].offset
		for k := i; k < j; k++ {
			out[k].offset = lo + (hi-lo)*float64(k-i+1)/float64(j-i+1)
		}
		i = j
	}
	for i := range out {
		out[i].offset = clampFloat(out[i].offset, 0, 1)
		if i > 0 && out[i].offset < out[i-1].offset {
			out[i].offset = out[i-1].offset
		}
	}
	return out
}
