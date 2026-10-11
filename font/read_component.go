package font

// readComponent reads one component of a composite glyph at p, returning it
// and the offset just past it. ok is false when the data is truncated.
func readComponent(data []byte, p int) (component, int, bool) {
	if p+4 > len(data) {
		return component{}, p, false
	}
	flags := be16(data, p)
	index := int(be16(data, p+2))
	p += 4

	var dx, dy float64
	if flags&0x0001 != 0 { // ARG_1_AND_2_ARE_WORDS
		if p+4 > len(data) {
			return component{}, p, false
		}
		dx, dy = float64(int16(be16(data, p))), float64(int16(be16(data, p+2)))
		p += 4
	} else {
		if p+2 > len(data) {
			return component{}, p, false
		}
		dx, dy = float64(int8(data[p])), float64(int8(data[p+1]))
		p += 2
	}
	if flags&0x0002 == 0 {
		// The arguments are point numbers to match up rather than an offset.
		// Almost nothing uses it and guessing would place the piece wrongly,
		// so it is placed at the origin instead.
		dx, dy = 0, 0
	}

	a, b, c, d, p, ok := readTransform(data, p, flags)
	if !ok {
		return component{}, p, false
	}
	return component{
		index: index, dx: dx, dy: dy,
		a: a, b: b, c: c, d: d,
		more: flags&0x0020 != 0,
	}, p, true
}
