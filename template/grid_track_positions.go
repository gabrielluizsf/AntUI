package template

import "github.com/gabrielluizsf/antui/template/css"

func gridTrackPositions(sizes []int, gap, total, available int, alignment uint8) []int {
	positions := make([]int, len(sizes))
	if len(sizes) == 0 {
		return positions
	}
	used := total + gap*(len(sizes)-1)
	left := max(available-used, 0)
	lead, between := 0, gap
	switch alignment {
	case css.ContentFlexEnd:
		lead = left
	case css.ContentCenter:
		lead = left / 2
	case css.ContentSpaceBetween:
		if len(sizes) > 1 {
			between += left / (len(sizes) - 1)
		}
	case css.ContentSpaceAround:
		between += left / len(sizes)
		lead = left / (2 * len(sizes))
	case css.ContentSpaceEvenly:
		between += left / (len(sizes) + 1)
		lead = left / (len(sizes) + 1)
	}
	position := lead
	for i := range sizes {
		positions[i] = position
		position += sizes[i] + between
	}
	return positions
}
