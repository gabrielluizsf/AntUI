package canvas

// maskRuns holds the working room a row needs, kept per depth: a shape, the
// clip it is cut by and the clip that one is cut by are never measured into
// the same row at once, and the depth they need is as many as the clips cut
// into clips along the way — which is never more than the picture is nested.
type maskRuns struct {
	box   int
	width int
	bufs  [][]float64
}

func (r *maskRuns) at(d int) []float64 {
	for len(r.bufs) <= d {
		r.bufs = append(r.bufs, make([]float64, r.width))
	}
	return r.bufs[d]
}

// cov is how much of one row of the band one measured shape covers: what its
// own outline covers, cut down to what every clip it carries covers too, one
// of them at a time. A pixel outside one of them multiplies out to nothing,
// which is the cut it is. It hands back the buffer it worked in at depth d,
// for the caller to take what it needs of before the next shape at that depth
// works in it again.
func (r *maskRuns) cov(m *maskMeasured, d, y int) []float64 {
	dst := r.at(d)
	clear(dst)
	for m.at < len(m.spans) && m.spans[m.at].y == y {
		s := m.spans[m.at]
		m.at++
		for x := s.x0; x < s.x1; x++ {
			if s.a > dst[x-r.box] {
				dst[x-r.box] = s.a
			}
		}
	}
	for j := range m.and {
		over := r.at(d + 1)
		clear(over)
		for k := range m.and[j] {
			mm := r.cov(&m.and[j][k], d+2, y)
			for x := range over {
				if mm[x] > over[x] {
					over[x] = mm[x]
				}
			}
		}
		for x := range dst {
			dst[x] *= over[x]
		}
	}
	return dst
}
