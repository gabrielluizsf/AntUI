package font

// simpleContours reads a glyph that is its own outline.
func simpleContours(data []byte, contours int) []Contour {
	if contours == 0 || 10+contours*2+2 > len(data) {
		return nil
	}
	ends := make([]int, contours)
	for i := range contours {
		ends[i] = int(be16(data, 10+i*2))
	}
	points := ends[contours-1] + 1
	if points <= 0 || points > 100000 {
		return nil
	}

	p := 10 + contours*2
	instructions := int(be16(data, p))
	p += 2 + instructions
	if p > len(data) {
		return nil
	}

	// Flags, which repeat: a flag with bit 3 set is followed by a count of
	// how many more points share it.
	flags := make([]byte, 0, points)
	for len(flags) < points {
		if p >= len(data) {
			return nil
		}
		flag := data[p]
		p++
		flags = append(flags, flag)
		if flag&0x08 != 0 {
			if p >= len(data) {
				return nil
			}
			repeat := int(data[p])
			p++
			for range repeat {
				if len(flags) >= points {
					break
				}
				flags = append(flags, flag)
			}
		}
	}

	// The coordinates are deltas, in one of three widths depending on the
	// flags: a byte, a byte that means "the same as last time", or a word.
	read := func(shortBit, sameBit byte) []float64 {
		out := make([]float64, points)
		v := 0.0
		for i, flag := range flags {
			switch {
			case flag&shortBit != 0:
				if p >= len(data) {
					return nil
				}
				d := float64(data[p])
				p++
				if flag&sameBit == 0 {
					d = -d
				}
				v += d
			case flag&sameBit == 0:
				if p+2 > len(data) {
					return nil
				}
				v += float64(int16(be16(data, p)))
				p += 2
			}
			out[i] = v
		}
		return out
	}
	xs := read(0x02, 0x10)
	ys := read(0x04, 0x20)
	if xs == nil || ys == nil {
		return nil
	}

	out := make([]Contour, 0, contours)
	start := 0
	for _, end := range ends {
		if end < start || end >= points {
			break
		}
		c := Contour{
			X:  xs[start : end+1],
			Y:  ys[start : end+1],
			On: make([]bool, end+1-start),
		}
		for i := range c.On {
			c.On[i] = flags[start+i]&0x01 != 0
		}
		out = append(out, c)
		start = end + 1
	}
	return out
}
