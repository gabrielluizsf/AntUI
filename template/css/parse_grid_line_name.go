package css

import (
	"strconv"
	"strings"
)

// parseGridLineName reads a line name and the span that may follow it, in any
// of the orders CSS allows: "sidebar-start", "-sidebar-start", "span 2",
// "2 span", "2 span 3", "sidebar-start span 2". The name keeps the case it was
// written with; only the span keyword folds.
func parseGridLineName(raw string) (GridLine, bool) {
	fields := strings.Fields(raw)
	if len(fields) == 0 || len(fields) > 3 {
		return GridLine{}, false
	}
	// The span keyword tells a number from a line to a count: a number written
	// before it names the line, a number written after it counts the tracks
	// spanned. A name is always the line, and a leading "-" counts that name
	// from the end of the axis.
	line := GridLine{Kind: GridLineAuto}
	spans, counts, count := 0, 0, 0
	for _, field := range fields {
		if strings.EqualFold(field, "span") {
			spans++
			continue
		}
		if n, err := strconv.Atoi(field); err == nil {
			if spans == 0 {
				if n == 0 || line.Kind != GridLineAuto {
					return GridLine{}, false
				}
				line = GridLine{Kind: GridLineIndex, Index: n}
				continue
			}
			if counts > 0 || n < 1 {
				return GridLine{}, false
			}
			counts, count = 1, n
			continue
		}
		if line.Kind != GridLineAuto {
			return GridLine{}, false
		}
		name, backward := gridBackwardName(field)
		if !validGridName(name) {
			return GridLine{}, false
		}
		line = GridLine{Kind: GridLineName, Name: name, Backward: backward}
	}
	if spans > 1 {
		return GridLine{}, false
	}
	if spans == 0 {
		// A lone number was read by parseGridLine itself, so only a name lives here.
		return line, line.Kind == GridLineName
	}
	line, ok := finishGridLine(line, count)
	return line, ok
}
