package css

import (
	"strconv"
	"strings"
)

// parseGridLine reads one <grid-line>: "auto", a line number counting from the start or the end, a line name, or any of them followed by a span.
// parseGridLine reads one <grid-line>: "auto", a line number counting from the
// start or the end, a line name, or any of them followed by a span.
func parseGridLine(raw string) (GridLine, bool) {
	raw = strings.TrimSpace(raw)
	if strings.EqualFold(raw, "auto") {
		return GridLine{Kind: GridLineAuto}, true
	}
	if n, err := strconv.Atoi(raw); err == nil {
		if n == 0 {
			return GridLine{}, false
		}
		return GridLine{Kind: GridLineIndex, Index: n}, true
	}
	return parseGridLineName(raw)
}
