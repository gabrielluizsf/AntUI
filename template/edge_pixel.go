package template

import "github.com/gabrielluizsf/antui/template/css"

// edgePixel is where one edge of a picture lands along an axis of its
// positioning area: percentages measure the side left over, fixed offsets the
// area itself, and a right or bottom anchor folds its offset back towards the
// inside, as the CSS "right 10px" inset does.
func (cs *cssStyle) edgePixel(st css.Style, pos css.BackPosition, u, box, img int) int {
	left := 0
	base := max(box-img, 0)
	switch pos.Edge {
	case css.BackPosMiddle:
		left = base / 2
	case css.BackPosEnd:
		left = base
	}
	if pos.Off.IsPct() {
		off := pos.Off.Resolve(css.Units{Width: base})
		if pos.Edge == css.BackPosEnd {
			left -= off
		} else {
			left += off
		}
	} else {
		off := pos.Off.Resolve(css.Units{Width: box}) * u
		if pos.Edge == css.BackPosEnd {
			left -= off
		} else {
			left += off
		}
	}
	return left
}
