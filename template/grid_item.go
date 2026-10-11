package template

import "github.com/gabrielluizsf/antui/template/css"

type gridItem struct {
	role, label string
	st          css.Style
	nw, nh      int
	contribW    int
	contribH    int
	contribMinW int
	contribMinH int
	w, h        int
	ml, mr      int
	mt, mb      int
	explicitW   bool
	explicitH   bool
	minW        int
	maxW        int
	minH        int
	maxH        int
	column      css.GridPlacement
	row         css.GridPlacement
	justifySelf uint8
	alignSelf   uint8
	order       int
	colStart    int
	colEnd      int
	rowStart    int
	rowEnd      int
	colDefinite bool
	rowDefinite bool
	out         bool
	inLeft      int
	inTop       int
	inRight     int
	inBottom    int
	x, y        int
	nested      *gridNested  // set when the item is a grid container of its own
	columns     *multiNested // set when the item is a column block of its own
	table       *tableNested // set when the item is a table of its own
}
