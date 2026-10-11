package svg

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
