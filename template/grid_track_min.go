package template

import "github.com/gabrielluizsf/antui/template/css"

// gridTrackMin is the floor a minmax gives its track. An auto floor is the
// automatic minimum size, which is the narrowest its items can be, and the
// content keywords are what they say they are.
func (b *gridBatch) gridTrackMin(size css.GridTrackSize, contribution gridContribution, base int, rowCyclic bool) int {
	switch size.Kind {
	case css.GridTrackLength:
		minimum := contribution.min
		if !rowCyclic || !size.Length.IsPct() {
			minimum = b.css.style.length(b.st, size.Length, base)
		}
		return max(minimum, 0)
	case css.GridTrackMinContent:
		return contribution.min
	case css.GridTrackMaxContent:
		return contribution.max
	default:
		return contribution.min
	}
}
