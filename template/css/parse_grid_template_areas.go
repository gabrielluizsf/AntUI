package css

import "strings"

// parseGridTemplateAreas reads the quoted row strings of grid-template-areas,
// validating every row's cell count and each named area's rectangular shape.
func parseGridTemplateAreas(raw string) ([][]string, bool) {
	raw = strings.TrimSpace(raw)
	if strings.EqualFold(raw, "none") {
		return nil, true
	}
	rows := splitGridList(raw, false)
	if len(rows) == 0 {
		return nil, false
	}
	out := make([][]string, 0, len(rows))
	width := -1
	for _, row := range rows {
		if len(row) < 2 || (row[0] != '"' && row[0] != '\'') || row[len(row)-1] != row[0] || !closedGridString(row) {
			return nil, false
		}
		cells := splitGridList(row[1:len(row)-1], false)
		if len(cells) == 0 {
			return nil, false
		}
		if width < 0 {
			width = len(cells)
		} else if len(cells) != width {
			return nil, false
		}
		for _, cell := range cells {
			if cell != "." && cell != "..." && !validGridIdent(cell) {
				return nil, false
			}
		}
		out = append(out, cells)
	}
	ok, _ := gridAreaBounds(out)
	return out, ok
}
