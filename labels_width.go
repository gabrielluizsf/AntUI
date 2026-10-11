package antui

// LabelsWidth is how wide a row of labels is: each with pad either side, and
// gap between one and the next.
func LabelsWidth(labels []string, pad, gap int) int {
	total := 0
	for i, label := range labels {
		if i > 0 {
			total += gap
		}
		total += LabelWidth(label, pad)
	}
	return total
}
