package template

import "github.com/gabrielluizsf/antui/template/css"

// measureTableContainer lays a table out on its own, so its parent has a box to
// place: the width the parent would stretch it to, and the height the columns
// and rows produce there. room is what the table has to fill, which a percentage
// of its own width answers to.
func (c *CSS) measureTableContainer(st css.Style, draw func(*CSS), room int) *tableNested {
	b := &tableBatch{css: c, st: st}
	c.setTableBox(b, st, 0, 0, room, 0, room)
	// The silent pass of the whole page is already running, so the rows are
	// collected here and the table is sized from them.
	c.tableWithin(b, draw)
	c.solveTable(b, 0, room, room)
	return &tableNested{batch: b, nw: c.tableAsk(b, room), nh: b.h}
}
