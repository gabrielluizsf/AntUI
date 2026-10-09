package css

import "strings"

// parseGridAutoFlow reads the grid-auto-flow flags: row/column plus the optional dense packing.
func parseGridAutoFlow(raw string) (uint8, bool, bool) {
	words := strings.Fields(raw)
	if len(words) == 0 || len(words) > 2 {
		return 0, false, false
	}
	var flow uint8
	var dense, haveFlow, haveDense bool
	for _, word := range words {
		switch word {
		case "row":
			if haveFlow {
				return 0, false, false
			}
			flow, haveFlow = GridAutoFlowRow, true
		case "column":
			if haveFlow {
				return 0, false, false
			}
			flow, haveFlow = GridAutoFlowColumn, true
		case "dense":
			if haveDense {
				return 0, false, false
			}
			dense, haveDense = true, true
		default:
			return 0, false, false
		}
	}
	return flow, dense, haveFlow
}
