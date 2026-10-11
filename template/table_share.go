package template

// tableShare is the width each column ends up with. The whole width is shared
// out in proportion to what the columns asked for, so a wide cell keeps more of
// the extra room than a narrow one, and the last column takes the rounding
// remainder so the columns add up to the width exactly.
func tableShare(cols []int, width int) []int {
	if len(cols) == 0 {
		return nil
	}
	natural := 0
	for _, w := range cols {
		natural += w
	}
	out := make([]int, len(cols))
	if natural == 0 {
		base, extra := width/len(cols), width%len(cols)
		for i := range out {
			out[i] = base
			if i < extra {
				out[i]++
			}
		}
		return out
	}
	shared := 0
	for i, w := range cols {
		if i == len(cols)-1 {
			out[i] = max(width-shared, 0)
			break
		}
		out[i] = width * w / natural
		shared += out[i]
	}
	return out
}
