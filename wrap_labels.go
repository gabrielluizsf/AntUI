package antui

// WrapLabels breaks a row of labels into rows none of which is wider than
// most, keeping the order they were given in. It gives back how many labels
// are on each row.
//
// A label too wide for a row of its own still gets one: making it share would
// only push something else out of the box as well, and a row that is too wide
// is a thing the caller can see and shorten. There is always at least one row
// when there is at least one label, so a caller can lay them out without
// checking for nothing.
func WrapLabels(labels []string, pad, gap, most int) []int {
	if len(labels) == 0 {
		return nil
	}
	rows := []int{0}
	width := 0
	for _, label := range labels {
		one := LabelWidth(label, pad)
		next := one
		if rows[len(rows)-1] > 0 {
			next = width + gap + one
		}
		if rows[len(rows)-1] > 0 && next > most {
			rows = append(rows, 0)
			next = one
		}
		rows[len(rows)-1]++
		width = next
	}
	return rows
}
