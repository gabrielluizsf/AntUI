package template

import "github.com/gabrielluizsf/antui/template/css"

type gridBatch struct {
	css      *CSS
	st       css.Style
	x, y     int
	w, h     int
	originX  int
	originY  int
	contentW int
	contentH int
	rowGap   int
	colGap   int
	items    []gridItem
	next     int
	// The explicit column tracks the last solve placed the items on, kept so a
	// container that has to size itself can ask how wide its own columns want
	// to be without resolving the template a second time.
	trackColumns []css.GridTrack
}
