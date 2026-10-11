package template

import "github.com/gabrielluizsf/antui/template/css"

// gridAreaName is the name a grid-area gave all four of its lines at once,
// which is how an item finds the template area it was pointed at.
func (it *gridItem) gridAreaName() string {
	if it.column.Start.Kind != css.GridLineName || it.column.End.Kind != css.GridLineName || it.row.Start.Kind != css.GridLineName || it.row.End.Kind != css.GridLineName {
		return ""
	}
	name := it.column.Start.Name
	if name == "" || it.column.End.Name != name || it.row.Start.Name != name || it.row.End.Name != name {
		return ""
	}
	return name
}
