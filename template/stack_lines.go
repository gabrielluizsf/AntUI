package template

import "github.com/gabrielluizsf/antui/template/css"

// stackLines answers where each line starts along the cross axis and how big
// the content grew: align-content spreads the leftover across the lines and
// stretch grows them to fill a definite cross size. The offsets are kept
// separate from the items so alignCross can add the within-line alignment.
func (b *flexBatch) stackLines(lines [][]int, crosses []int, gap, avail int, definite bool) ([]int, int) {
	if len(lines) == 0 {
		return nil, 0
	}
	total := 0
	for i, c := range crosses {
		if i > 0 {
			total += gap
		}
		total += c
	}
	if definite && b.st.AlignContent == css.ContentStretch {
		// Stretch fills the cross size the container was given: the lines are
		// made to add up to it, so a line wider than the container is narrowed
		// to fit and what does not fit spills out of it rather than the
		// container growing to hold it.
		room := max(avail-gap*(len(lines)-1), 0) / len(lines)
		for i := range crosses {
			crosses[i] = room
		}
		total = 0
		for i, c := range crosses {
			if i > 0 {
				total += gap
			}
			total += c
		}
	}
	left := 0
	if definite {
		left = max(avail-total, 0)
	}
	lead, between := 0, gap
	switch b.st.AlignContent {
	case css.ContentFlexEnd:
		lead = left
	case css.ContentCenter:
		lead = left / 2
	case css.ContentSpaceBetween:
		if len(lines) > 1 {
			between = gap + left/(len(lines)-1)
		}
	case css.ContentSpaceAround:
		between = gap + left/len(lines)
		lead = left / (2 * len(lines))
	case css.ContentSpaceEvenly:
		between = gap + left/(len(lines)+1)
		lead = left / (len(lines) + 1)
	}
	offsets := make([]int, len(lines))
	pos := lead
	for li := range lines {
		offsets[li] = pos
		if li < len(lines)-1 {
			pos += crosses[li] + between
		}
	}
	if definite {
		return offsets, avail
	}
	return offsets, total
}
