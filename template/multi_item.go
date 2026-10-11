package template

import "github.com/gabrielluizsf/antui/template/css"

// multiItem is one child of a multi-column container as the flow sees it: the
// box it asked for, and the pieces the flow cut that box into. A piece is one
// run of the box in one column, and a box the flow did not have to cut keeps
// exactly one.
type multiItem struct {
	role, label string
	st          css.Style
	nw, nh      int // measured border-box size, untouched by the flow
	w, h        int // solved: the width the column gives, the height left
	ml, mr      int // margins across the inline axis
	mt, mb      int // margins along the block axis
	keep        bool
	pieces      []multiPiece
	next        int
	nested      *multiNested // set when the child is a container of its own
	table       *tableNested // set when the child is a table
}
