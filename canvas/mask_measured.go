package canvas

// maskMeasured is one shape measured over the whole band, along with the clips
// it carries and the clips those carry, all the way down. The runs of one list
// come back in row order, so each is walked from where it left off rather than
// being searched for again on every row.
type maskMeasured struct {
	at    int
	spans []span
	and   [][]maskMeasured
}

func measureShape(s MaskShape, bx0, by0, bx1, by1 int) maskMeasured {
	w := maskMeasured{spans: fillCoverage(s.Path, bx0, by0, bx1, by1, s.Rule)}
	for _, e := range s.And {
		entry := make([]maskMeasured, 0, len(e))
		for _, es := range e {
			if es.Path == nil {
				continue
			}
			entry = append(entry, measureShape(es, bx0, by0, bx1, by1))
		}
		w.and = append(w.and, entry)
	}
	return w
}
