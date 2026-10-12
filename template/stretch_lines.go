package template

// stretchLines makes the lines add up to the definite cross size, widening a
// line that has room and narrowing one that does not, and returns the new
// total.
func stretchLines(crosses []int, gap, avail int) int {
	n := len(crosses)
	room := max(avail-gap*(n-1), 0) / n
	for i := range crosses {
		crosses[i] = room
	}
	total := 0
	for i, c := range crosses {
		if i > 0 {
			total += gap
		}
		total += c
	}
	return total
}
