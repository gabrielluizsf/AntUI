package template

import "github.com/gabrielluizsf/antui/template/css"

// multiBatch is one multi-column container: its box on the page, the columns
// the flow made, and the items collected while its children drew. The silent
// pass measures every child, the flow cuts the boxes into pieces, and the draw
// passes hand those pieces back — one pass per fragment an item was cut into,
// so every widget is drawn once per piece it owns.
type multiBatch struct {
	css      *CSS
	st       css.Style
	x, y, w  int // container border box
	h        int // solved after the columns have a height
	originX  int // content-box top-left in window coordinates
	originY  int
	contentW int
	contentH int
	left     int // the border and padding beside the columns
	right    int
	gap      int  // the gutter between two columns
	count    int  // how many columns there are
	colW     int  // the width of each one
	top      int  // the border and padding above the columns
	bottom   int  // the border and padding below them
	room     int  // the height one column may hold
	colH     int  // the tallest column's content
	cuts     bool // a height decided outside the flow: a box that does not fit is cut
	definite bool // the height of the box was decided, not asked for
	items    []multiItem
	next     int // the item the draw pass is handing back
	pass     int // which fragment pass the draw is on
	passes   int // how many passes the pieces need
}
