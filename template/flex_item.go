package template

import "github.com/gabrielluizsf/antui/template/css"

// flexItem is one child of a flex container as the solver sees it: the box
// the stylesheet wants, collected in the silent pass and turned into a final
// x,y,w,h in the draw one. Every field the layout mutates is reset fresh by
// solve from nw/nh, so the two passes always answer the same geometry.
type flexItem struct {
	role, label  string
	st           css.Style
	nw, nh       int // measured border-box size, untouched by the solver
	w, h         int // solved border-box size
	ml, mr       int // margins on the main axis
	mt, mb       int // margins on the cross axis
	order        int
	grow, shrink float64
	basis        css.Length
	alignSelf    uint8 // AlignAuto means follow the container
	explicitW    bool  // width declared: stretch must not widen it
	explicitH    bool  // height declared: stretch must not heighten it
	minMain      int   // explicit main-axis minimum, -1 when absent
	maxMain      int   // explicit main-axis maximum, -1 when absent
	base         int   // main-axis size before grow or shrink
	out          bool  // out of flow: positioned by insets, not by the solver
	inLeft       int   // out-of-flow insets, -1 when absent
	inTop        int
	inRight      int
	inBottom     int
	x, y         int          // solved position relative to the content origin
	nested       *gridNested  // set when the item is a grid container of its own
	columns      *multiNested // set when the item is a column block of its own
	table        *tableNested // set when the item is a table of its own
}
