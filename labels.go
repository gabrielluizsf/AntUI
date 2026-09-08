package antui

import "github.com/gabrielluizsf/antui/canvas"

// Measuring a row of labels, which is what sizing a thing to what is in it
// comes down to.
//
// A dialog, a menu, a toolbar and a set of tabs are all the same shape: some
// pieces of text, each with room either side, laid along a row with a gap
// between them. Every one of them has to answer two questions — how wide is
// that, and what happens when it will not fit — and every one of them used to
// answer them with a guess. A guess is fine until somebody puts three answers
// on a dialog and the last one is drawn outside the box.
//
// The measuring lives here because this is where the font is. Nothing above
// it can measure text without asking, so nothing above it should have to
// invent a number.

// LabelWidth is how wide a piece of text is with pad either side of it.
func LabelWidth(text string, pad int) int {
	return canvas.TextWidth(text) + pad*2
}

// LabelHeight is how tall one is with pad above and below.
func LabelHeight(pad int) int { return canvas.FontHeight + pad*2 }

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
