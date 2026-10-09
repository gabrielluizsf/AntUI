package css

import "strings"

// parseGridArea reads the grid-area shorthand into its row and column placements.
func parseGridArea(raw string) (GridPlacement, GridPlacement, bool) {
	raw = strings.TrimSpace(raw)
	if strings.EqualFold(raw, "none") || strings.EqualFold(raw, "auto") {
		return GridPlacement{}, GridPlacement{}, true
	}
	parts := splitGridSlashes(raw)
	if len(parts) == 0 || len(parts) > 4 {
		return GridPlacement{}, GridPlacement{}, false
	}
	lines := make([]GridLine, len(parts))
	for i, part := range parts {
		line, ok := parseGridLine(part)
		if !ok {
			return GridPlacement{}, GridPlacement{}, false
		}
		lines[i] = line
	}
	if len(lines) == 1 {
		area := GridPlacement{Start: lines[0]}
		if lines[0].Kind == GridLineName {
			// A lone name is the name of a template area: it stands for all
			// four lines at once.
			area.End = lines[0]
			return area, area, true
		}
		return GridPlacement{}, area, true
	}
	row := GridPlacement{Start: lines[0]}
	column := GridPlacement{Start: lines[1]}
	if len(lines) >= 3 {
		row.End = lines[2]
	}
	if len(lines) == 4 {
		column.End = lines[3]
	}
	return column, row, true
}
