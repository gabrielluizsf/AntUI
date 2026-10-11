package svg

import "github.com/gabrielluizsf/antui/canvas"

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
