package canvas

// dasher walks one run of a dashed path, opening and closing stretches of it
// as the pattern turns on and off. on is whether the stretch being walked is
// drawn, open whether the run of it has been opened yet; every even phase is a
// dash.
type dasher struct {
	p       *Path
	pattern []float64
	phase   int
	rest    float64
	on      bool
	open    bool
}

// cut walks one segment of the run, cutting it at every pattern stretch it
// holds and drawing the ones that are on.
func (d *dasher) cut(a, b Point, seg float64) {
	walked := 0.0
	// A segment long enough to hold a whole number of pattern stretches is cut
	// at each of them; a segment shorter than what is left of the current one
	// is only walked to its end.
	for seg-walked > d.rest {
		part := d.rest + walked
		at := Point{a.X + (b.X-a.X)*part/seg, a.Y + (b.Y-a.Y)*part/seg}
		if d.on {
			if d.open {
				d.p.LineTo(at.X, at.Y)
			} else {
				d.p.MoveTo(at.X, at.Y)
			}
		}
		d.phase = (d.phase + 1) % len(d.pattern)
		d.rest = d.pattern[d.phase]
		d.on = d.phase%2 == 0
		// A dash opens wherever the pattern comes back on, which is either the
		// point the last one stopped at or the start of the line.
		if d.on {
			d.p.MoveTo(at.X, at.Y)
		}
		d.open = d.on
		walked = part
	}
	d.rest -= seg - walked
	if d.on {
		d.p.LineTo(b.X, b.Y)
	} else {
		d.open = false
	}
}
