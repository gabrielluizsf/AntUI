package canvas

import "math"

// Dash is a stroke broken into pieces along its length, the dashes and the gaps
// between them measured in the same units as the path. A pattern of one number
// alternates that length on and the same length off; two numbers are a dash and
// a gap; an odd list is read twice over, which is how "ten on, five off" becomes
// dashes of ten with gaps of five and then gaps of ten with dashes of five. A
// list that is empty, or that holds a zero or a negative, is not a dash at all
// and the line comes back whole.
type Dash struct {
	On, Off []float64
	Offset  float64

	// The pattern as one run of dashes and gaps, and the distance a whole turn
	// of it measures, worked out once so that the walk does not build it again
	// for every point of every subpath.
	pattern []float64
	period  float64
}

// NewDash makes a dash pattern out of the lengths stroke-dasharray was given,
// and reports whether it is a pattern at all. An odd list is doubled rather
// than refused, which is what lets a single number mean "dashes as long as the
// gaps".
func NewDash(lengths []float64, offset float64) (Dash, bool) {
	if len(lengths) == 0 {
		return Dash{}, false
	}
	pattern := append([]float64{}, lengths...)
	if len(pattern)%2 == 1 {
		pattern = append(pattern, pattern...)
	}
	// The period is what the pattern as it will be walked measures, so it is
	// counted after an odd list has been doubled: "ten on, five off" repeated
	// comes to twenty-five, not to fifteen.
	sum := 0.0
	for _, l := range pattern {
		if l <= 0 || math.IsNaN(l) {
			return Dash{}, false
		}
		sum += l
	}
	if sum <= 0 {
		return Dash{}, false
	}
	// An odd pattern is read twice over, so the dashes and the gaps are the two
	// halves of it and there is one of each of them.
	halves := len(pattern) / 2
	d := Dash{
		On:     append([]float64{}, pattern[:halves]...),
		Off:    append([]float64{}, pattern[halves:]...),
		Offset: wrap(offset, sum),
	}
	d.pattern, d.period = pattern, sum
	return d, true
}

// pattern is the dash pattern as one run, dashes and gaps in turn, and period
// how far it has to go before it starts again. Both are worked out once, because
// every step of the walk asks for them.
func (d Dash) cycle() ([]float64, float64) {
	if len(d.pattern) > 0 {
		return d.pattern, d.period
	}
	// A pattern built by hand rather than by NewDash still has to be read.
	out := make([]float64, 0, len(d.On)*2)
	for i, l := range d.On {
		out = append(out, l, d.Off[i%len(d.Off)])
	}
	sum := 0.0
	for _, l := range out {
		sum += l
	}
	return out, sum
}

// DashedPath is the path a dashed stroke draws: the same line cut into the runs
// that are on, each as a subpath of its own so that the two ends of a dash can
// be capped. A line that is not dashed comes back as it was, and a path of
// nothing is still a path of nothing.
func (p *Path) Dashed(d Dash) *Path {
	if p == nil {
		return nil
	}
	pattern, period := d.cycle()
	if len(pattern) == 0 || period <= 0 {
		return p
	}
	pts, closed := p.Points()
	out := NewPath()
	for i, sub := range pts {
		if len(sub) < 2 {
			if len(sub) == 1 {
				out.MoveTo(sub[0].X, sub[0].Y)
			}
			continue
		}
		if closed[i] {
			sub = append(append([]Point{}, sub...), sub[0])
		}
		// The pattern is walked from the start for every subpath, so a dash
		// that would have begun a while ago begins a while ago here too —
		// which is what a shape drawn as one path expects and what treating
		// each subpath as its own line would get wrong.
		out.appendDashes(sub, pattern, d.Offset, period)
	}
	return out
}

// appendDashes cuts one run of points into the stretches of it that are on,
// beginning at whatever point of the pattern the offset lands on. The walk goes
// by distance along the run rather than by its points, because a dash falls
// where the distances add up, and the points of a curve are not evenly spaced.
func (p *Path) appendDashes(sub []Point, pattern []float64, offset, period float64) {
	// Which stretch of the pattern the line starts in, and how much of it is
	// already gone by the time the first point is reached.
	phase, gone := 0, wrap(offset, period)
	for {
		phase %= len(pattern)
		if gone < pattern[phase] {
			break
		}
		gone -= pattern[phase]
		phase++
	}
	rest := pattern[phase] - gone

	// The offset can land past whole segments of the line, so the walk begins
	// where it lands and the stretch before that point is on no dash at all.
	from := 0
	for from+1 < len(sub) {
		seg := math.Hypot(sub[from+1].X-sub[from].X, sub[from+1].Y-sub[from].Y)
		if gone < seg {
			break
		}
		gone, from = gone-seg, from+1
	}
	if from+1 >= len(sub) {
		return
	}
	at := sub[from]
	if gone > 0 {
		a, b := sub[from], sub[from+1]
		seg := math.Hypot(b.X-a.X, b.Y-a.Y)
		at = Point{a.X + (b.X-a.X)*gone/seg, a.Y + (b.Y-a.Y)*gone/seg}
	}

	// Whether the stretch being walked is one that is drawn, and whether the run
	// of it has been opened yet. Every even phase is a dash.
	on, open := phase%2 == 0, false
	if on {
		p.MoveTo(at.X, at.Y)
		open = true
	}
	for i := from + 1; i < len(sub); i++ {
		a, b := sub[i-1], sub[i]
		seg := math.Hypot(b.X-a.X, b.Y-a.Y)
		// The segment the offset landed in is cut short by however far into it
		// the walk had already come.
		if i == from+1 {
			a = at
			seg -= gone
		}
		walked := 0.0
		// A segment long enough to hold a whole number of pattern stretches
		// is cut at each of them; a segment shorter than what is left of the
		// current one is only walked to its end.
		for seg-walked > rest {
			part := rest + walked
			at := Point{a.X + (b.X-a.X)*part/seg, a.Y + (b.Y-a.Y)*part/seg}
			if on {
				if open {
					p.LineTo(at.X, at.Y)
				} else {
					p.MoveTo(at.X, at.Y)
				}
			}
			phase = (phase + 1) % len(pattern)
			rest = pattern[phase]
			on = phase%2 == 0
			// A dash opens wherever the pattern comes back on, which is
			// either the point the last one stopped at or the start of the line.
			if on {
				p.MoveTo(at.X, at.Y)
			}
			open = on
			walked = part
		}
		rest -= seg - walked
		if on {
			p.LineTo(b.X, b.Y)
		} else {
			open = false
		}
	}
}

// wrap is a distance along a pattern brought back inside it, counting from the
// start even when the distance is given as a negative one: an offset of minus
// half a period is where half a period before the end is, not before the start.
func wrap(d, period float64) float64 {
	if d = math.Mod(d, period); d < 0 {
		d += period
	}
	return d
}
