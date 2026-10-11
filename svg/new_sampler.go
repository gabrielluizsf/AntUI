package svg

import "github.com/gabrielluizsf/antui/canvas"

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
