package canvas

// normalizeGradientStops fills the offsets a caller left out the way CSS does:
// the first and last missing stops become 0 and 1, and missing stops between
// known ones spread evenly over the gap.
func normalizeGradientStops(in []GradientStop) []GradientStop {
	n := len(in)
	if n == 0 {
		return nil
	}
	if n == 1 {
		return []GradientStop{{Offset: 0, Color: in[0].Color}, {Offset: 1, Color: in[0].Color}}
	}
	out := make([]GradientStop, n)
	copy(out, in)
	has := make([]bool, n)
	any := false
	for i := range n {
		has[i] = out[i].Offset >= 0
		any = any || has[i]
	}
	if !any {
		for i := range n {
			out[i].Offset = float64(i) / float64(n-1)
		}
		return out
	}
	if !has[0] {
		out[0].Offset = 0
		has[0] = true
	}
	if !has[n-1] {
		out[n-1].Offset = 1
		has[n-1] = true
	}
	// Fill each run of unset stops between known neighbours.
	i := 0
	for i < n {
		if has[i] {
			i++
			continue
		}
		j := i
		for j < n && !has[j] {
			j++
		}
		lo, hi := out[i-1].Offset, out[j].Offset
		for k := i; k < j; k++ {
			f := float64(k-i+1) / float64(j-i+1)
			out[k].Offset = lo + (hi-lo)*f
		}
		i = j
	}
	return out
}
