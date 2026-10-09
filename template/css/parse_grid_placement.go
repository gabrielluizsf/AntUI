package css

// parseGridPlacement reads grid-row/column: a start line and an optional end line.
func parseGridPlacement(raw string) (GridPlacement, bool) {
	parts := splitGridSlashes(raw)
	if len(parts) == 0 || len(parts) > 2 {
		return GridPlacement{}, false
	}
	placement := GridPlacement{}
	var ok bool
	if placement.Start, ok = parseGridLine(parts[0]); !ok {
		return GridPlacement{}, false
	}
	if len(parts) == 2 {
		if placement.End, ok = parseGridLine(parts[1]); !ok {
			return GridPlacement{}, false
		}
	}
	return placement, true
}

// parseGridArea reads a grid-area: a single value is one line — a name standing
// for the item's four lines, or a row line; two to four values are line
// placements written row first — "row-start / column-start [ / row-end /
// column-end ]".
