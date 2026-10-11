package template

import "github.com/gabrielluizsf/antui/template/css"

// gridAxisPlacement reads a placement off one axis into the range of tracks it
// covers. A line is an index counted from the start or the end, or a name the
// template carries — the first line carrying the name for a start, the last for
// an end, which is how a name bounds both sides of an item. A name no line
// carries is no placement at all, leaving the item to be placed automatically.
func gridAxisPlacement(axis css.GridAxis, placement css.GridPlacement) (int, int, bool) {
	start, startOK := gridAxisLine(axis, placement.Start, false)
	end, endOK := gridAxisLine(axis, placement.End, true)
	if !startOK && !endOK {
		return 0, 0, false
	}

	if !startOK {
		start = end - gridPlacementSpan(placement)
	}
	if !endOK {
		if placement.End.Kind == css.GridLineSpan {
			end = start + placement.End.Index
		} else if placement.Start.Span > 0 {
			end = start + placement.Start.Span
		} else {
			end = start + 1
		}
	}
	if end < start {
		start, end = end, start
	}
	if end <= start {
		end = start + 1
	}
	return start, end, true
}
