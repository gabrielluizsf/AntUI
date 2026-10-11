package template

import "github.com/gabrielluizsf/antui/template/css"

// gridAxisLine is the line one side of a placement names, and whether the
// placement named one at all. A count is taken from the start unless it is
// negative, and a name is looked up from the end of the axis when the end of
// the item is what named it or when a "-" sent it there.
func gridAxisLine(axis css.GridAxis, line css.GridLine, end bool) (int, bool) {
	explicit := len(axis.Tracks)
	switch line.Kind {
	case css.GridLineIndex:
		if line.Index > 0 {
			return line.Index - 1, true
		}
		if line.Index < 0 {
			return explicit + line.Index + 1, true
		}
	case css.GridLineName:
		if found, ok := axis.LineIndex(line.Name, end || line.Backward); ok {
			return found, true
		}
	}
	return 0, false
}
