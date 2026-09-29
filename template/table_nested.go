package template

import "github.com/gabrielluizsf/antui/template/css"

// A table drawn inside another container is an item of it, not a flattening of
// its cells into the parent: the parent's solver picks the box the table lands
// in and the table lines its own rows up inside it. A table inside a cell is
// that cell's content, the same way a column block is.
//
// The silent pass measures the container where the parent would stretch it —
// the width a percentage inside it answers to — and the draw pass re-solves it
// at the box the parent really gave it.

// tableNested is one table drawn inside another container. It carries the batch
// that lines its rows up — collected by the silent pass and re-solved at the
// final box — and the border box the table asks its parent for.
type tableNested struct {
	batch  *tableBatch
	nw, nh int
}

// nestedTable is [CSS.Table] called from inside another container: the table
// becomes an item of the parent, placed like any other and holding its own
// rows.
func (c *CSS) nestedTable(draw func(*CSS)) {
	st := c.style.baseStyle(css.RoleTable, State{})
	if st.Display == css.DisplayNone {
		return
	}
	switch {
	case c.multicoling:
		c.tableInColumns(st, draw)
	case c.flexing:
		c.flexNestedTable(st, draw)
	case c.griding:
		c.gridNestedTable(st, draw)
	case c.row != nil:
		c.tableInCell(st, draw)
	default:
		// Straight in a table's callback: a row of its own, holding the table.
		c.tableRowOfOneBlock(func(*CSS) { c.tableInCell(st, draw) })
	}
}

// tableRowOfOneBlock wraps a block drawn straight in a table in the row it
// belongs to — a row of a single cell, which is what a browser's anonymous
// boxes make of a child that did not ask to be a row. collect runs with the row
// in place, so the block is measured in the silent pass and placed in the draw
// pass like any other cell.
func (c *CSS) tableRowOfOneBlock(collect func(*CSS)) {
	b := c.table
	if b == nil {
		return
	}
	if c.measure {
		r := &tableRow{}
		c.rowWithin(b, r, collect)
		b.children = append(b.children, tableChild{row: r})
		return
	}
	if b.rowIdx >= len(b.children) {
		return
	}
	child := b.children[b.rowIdx]
	b.rowIdx++
	if child.row == nil {
		return
	}
	c.rowWithin(b, child.row, collect)
}

// tableInCell puts a table inside the cell of the row in place. The cell is
// the table's own box, so the row gives it the width of its column and the
// height the table asks for.
func (c *CSS) tableInCell(st css.Style, draw func(*CSS)) {
	if c.measure {
		nested := c.measureTableContainer(st, draw, c.cellRoom())
		cell := c.tableCellSized(st, nested.nw, nested.nh)
		cell.table = nested
		c.row.cells = append(c.row.cells, cell)
		return
	}
	if c.row.next >= len(c.row.cells) {
		return
	}
	cell := c.row.cells[c.row.next]
	c.row.next++
	if cell.table == nil {
		return
	}
	b := c.table
	// A table is as wide as its own columns ask for, not as wide as the cell, so
	// it centres in the cell rather than filling it.
	c.drawTableContainer(cell.table, draw, cell.w, b.originX+cell.x, b.originY+cell.y, cell.h)
}

// multiInTableCell puts a column block inside the cell of the row in place. A
// block drawn straight in the table's callback has no row to land in, so it is
// wrapped in a row of its own first, like a table that founds the same way.
func (c *CSS) multiInTableCell(st css.Style, draw func(*CSS)) {
	if c.row == nil {
		c.tableRowOfOneBlock(func(*CSS) { c.multiInTableCell(st, draw) })
		return
	}
	if c.measure {
		nested := c.measureMultiContainer(st, draw, c.cellRoom())
		cell := c.tableCellSized(st, nested.nw, nested.nh)
		cell.multi = nested
		c.row.cells = append(c.row.cells, cell)
		return
	}
	if c.row.next >= len(c.row.cells) {
		return
	}
	cell := c.row.cells[c.row.next]
	c.row.next++
	if cell.multi == nil {
		return
	}
	b := c.table
	x, y, w, h := b.originX+cell.x, b.originY+cell.y, cell.w, cell.h
	c.drawMultiContainer(cell.multi, draw, w, x, y, w, h)
}

// cellRoom is the width a percentage inside a cell answers to: the space the
// row has for one column of the table, or the room the flow has for the table
// when its columns are not solved yet.
func (c *CSS) cellRoom() int {
	b := c.table
	if b == nil {
		return c.win.Width()
	}
	if len(b.cols) == 0 {
		return max(b.contentW, c.win.Width())
	}
	return b.cols[min(len(c.row.cells), len(b.cols)-1)]
}

// tableCellSized is a cell whose content is a container: it asks for the box
// the container asked for, and the row gives it the width of its column.
func (c *CSS) tableCellSized(st css.Style, nw, nh int) *tableCell {
	return &tableCell{st: st, nw: nw, nh: nh, w: nw, h: nh}
}

func (c *CSS) flexNestedTable(st css.Style, draw func(*CSS)) {
	parent := c.flex
	if c.measure {
		nested := c.measureTableContainer(st, draw, parent.contentW)
		it := c.flexItemSized(st, nested.nw, nested.nh)
		it.table = nested
		parent.items = append(parent.items, it)
		return
	}
	if parent.next >= len(parent.items) {
		return
	}
	it := &parent.items[parent.next]
	if it.table == nil {
		return
	}
	parent.next++
	c.drawTableContainer(it.table, draw, it.w, parent.originX+it.x, parent.originY+it.y, it.h)
}

func (c *CSS) gridNestedTable(st css.Style, draw func(*CSS)) {
	parent := c.grid
	if c.measure {
		nested := c.measureTableContainer(st, draw, parent.contentW)
		it := c.gridItemSized("", "", st, nested.nw, nested.nh, nested.nw, nested.nh)
		it.table = nested
		parent.items = append(parent.items, it)
		return
	}
	if parent.next >= len(parent.items) {
		return
	}
	it := &parent.items[parent.next]
	if it.table == nil {
		return
	}
	parent.next++
	c.drawTableContainer(it.table, draw, it.w, parent.originX+it.x, parent.originY+it.y, it.h)
}

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

// tableAsk is how wide a table being measured says it wants to be to whatever
// holds it. A width of its own answers for itself; a width in a percentage is a
// fraction of a box the parent has not solved yet — the room the measure pass
// offered is a guess, and in a cell whose columns are not sized it is the whole
// window. A percentage with no definite box to resolve against counts as auto,
// so the answer is the width the columns in the table ask for, which is what
// sizes the column that holds it. The draw pass solves the table again at the
// width the parent really gave, and the percentage answers to that.
func (c *CSS) tableAsk(b *tableBatch, room int) int {
	if !b.st.Has("width") || !b.st.Width.IsPct() {
		return b.w
	}
	left, right, _, _ := c.tableChrome(b.st, room)
	return colsWidth(c.tableColumns(b)) + left + right
}

// drawTableContainer places a nested table in the room its parent solved for it
// — the cell, column or item it was drawn in — solves it again so the columns
// answer to that room, paints its face and lines its rows up inside.
func (c *CSS) drawTableContainer(nested *tableNested, draw func(*CSS), room, x, y, h int) {
	b := nested.batch
	b.y = y
	c.setTableBox(b, b.st, x, y, room, 0, room)
	c.solveTable(b, x, room, room)
	if c.style.boxOn(b.st) && visible(b.st) {
		c.style.paintBox(c.win, b.st, b.x, b.y, b.w, b.h)
	}
	c.tableWithin(b, draw)
}
