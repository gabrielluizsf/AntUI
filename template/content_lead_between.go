package template

import "github.com/gabrielluizsf/antui/template/css"

// contentLeadBetween turns the leftover cross space into the lead before the
// first line and the gap between the lines, following align-content.
func contentLeadBetween(align uint8, left, n, gap int) (int, int) {
	lead, between := 0, gap
	switch align {
	case css.ContentFlexEnd:
		lead = left
	case css.ContentCenter:
		lead = left / 2
	case css.ContentSpaceBetween:
		if n > 1 {
			between = gap + left/(n-1)
		}
	case css.ContentSpaceAround:
		between = gap + left/n
		lead = left / (2 * n)
	case css.ContentSpaceEvenly:
		between = gap + left/(n+1)
		lead = left / (n + 1)
	}
	return lead, between
}
