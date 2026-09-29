package template

import "github.com/gabrielluizsf/antui/template/css"

// A table is a block in the flow that widens to its content: the widgets its
// callback draws are cells, and the cells a row draws are lined up in columns
// as wide as the widest cell in their column, the way a browser lays out a
// table with separated borders. The rows come out as tall as their tallest
// cell, the columns as wide as the whole table the stylesheet asked for, and a
// cell's content is placed in its box by the vertical-align the cell's own
// style declares.
//
//	app.Table(func(t *CSS) {
//		for _, row := range rows {
//			t.TableRow(func(r *CSS) {
//				for _, cell := range row {
//					r.Label(cell)
//				}
//			})
//		}
//	})
//
// A row the stylesheet marks as a row — display: table-row on the widget — is
// laid out as one; anything else drawn directly in the table becomes a row of
// its own, which is what a browser's anonymous boxes do. A widget that says
// display: table-caption is not a cell: it is painted above the table, across
// its whole width.

// tableCell is one cell of a row: the box it wants, the box the row gave it,
// and the container it holds when the cell's content is a table or a column
// block of its own.
type tableCell struct {
	role, label string
	st          css.Style
	nw, nh      int // measured border-box size
	w, h        int // solved by the row
	x, y        int // solved, relative to the table's content origin
	multi       *multiNested
	table       *tableNested
}

// tableRow is one row: the cells drawn inside it, and the height the tallest
// one gives it.
type tableRow struct {
	cells  []*tableCell
	height int
	next   int // the cell the draw pass is handing back
}

// tableChild is one thing a table draws between its rows: a row, or the
// caption the table paints above itself.
type tableChild struct {
	row     *tableRow
	caption *tableCell
}

// tableBatch is one table: its box, the columns the solver cut out of it and
// the rows the callback drew. The silent pass collects the cells, the solver
// gives every cell the width of its column and every row the height of its
// tallest cell, and the draw pass hands the boxes back.
type tableBatch struct {
	css      *CSS
	st       css.Style
	x, y, w  int // container border box
	h        int
	originX  int // content-box top-left in window coordinates
	originY  int
	top      int // the border and padding above the rows
	bottom   int // the border and padding below them
	caption  int // the height a caption takes above the rows
	contentW int
	contentH int
	cols     []int // the width each column was given
	children []tableChild
	row      *tableRow // the row collecting cells, nil outside one
	rowIdx   int       // the row the draw pass is on
}

// Table lays the widgets its callback draws out as a table: a block in the
// flow whose rows are drawn by [CSS.TableRow], each cell the widget drawn for
// it, the columns as wide as the widest cell they hold and the rows as tall as
// the tallest cell in them. The table is styled by the table role, and the
// cells by the roles the widgets drawn inside it carry.
func (c *CSS) Table(draw func(*CSS)) {
	if c.multicoling || c.flexing || c.griding || c.tabling {
		c.nestedTable(draw)
		return
	}

	st := c.style.baseStyle(css.RoleTable, State{})
	if st.Display == css.DisplayNone {
		return
	}
	winW := c.win.Width()
	mt := c.style.length(st, st.Margin[0], winW)
	mb := c.style.length(st, st.Margin[2], winW)

	// The table is a block: it closes the inline line, collapses the sibling
	// margins and rests on the flow cursor.
	c.cursorY += c.lineH
	c.lineH = 0
	c.lineX = 0
	c.cursorY += max(c.lastMB, mt)
	y := c.cursorY

	b := &tableBatch{css: c, st: st}
	c.setTableBox(b, st, 0, y, 0, 0, winW)
	c.table, c.tabling = b, true

	c.measure = true
	draw(c)
	c.measure = false

	c.solveTable(b, 0, winW, winW)

	if c.style.boxOn(st) && visible(st) {
		c.style.paintBox(c.win, st, b.x, b.y, b.w, b.h)
	}
	c.tableWithin(b, draw)

	c.tabling = false
	c.table = nil
	c.cursorY = b.y + b.h
	c.lastMB = mb
}

// setTableBox gives a batch the geometry of the table's border box: where it
// sits, how wide it is and the content box the cells are lined up in. The
// caption, which a table paints above itself, is part of that content box, so
// the rows start below it. pctBase is what a percentage of the table answers to.
func (c *CSS) setTableBox(b *tableBatch, st css.Style, x, y, w, h, pctBase int) {
	left, right, top, bottom := c.tableChrome(st, pctBase)
	b.x, b.y, b.w, b.h = x, y, w, h
	b.top, b.bottom = top, bottom
	b.originX, b.originY = x+left, y+top+b.caption
	b.contentW = max(w-left-right, 0)
	b.contentH = max(h-top-bottom-b.caption, 0)
}

// tableWithin runs the children of one table: the callback draws them, the
// silent pass collects the rows and the draw pass hands the cells their boxes.
// The batch in place is the table they belong to, so a table inside this one is
// collected as a cell of its own. The pass the page is in carries through, as it
// does in multiWithin: a nested table reached by the silent pass collects its own
// rows there, and one reached by a draw pass places them.
func (c *CSS) tableWithin(b *tableBatch, draw func(*CSS)) {
	outer, outerTable, outerRow, outerMulti := c.table, c.tabling, c.row, c.multicoling
	outerFlex, outerGrid := c.flexing, c.griding
	c.table, c.tabling, c.row, c.multicoling = b, true, nil, false
	c.flexing, c.griding = false, false
	b.rowIdx = 0
	draw(c)
	c.table, c.tabling, c.row, c.multicoling = outer, outerTable, outerRow, outerMulti
	c.flexing, c.griding = outerFlex, outerGrid
}

// TableRow draws one row of a table: the widgets its callback draws are the
// cells of the row, lined up in the columns they belong to. A row is a block of
// its own — the widgets before it and after it are other rows — and the row is
// as tall as the tallest cell in it.
//
// The callback runs in the silent pass to measure the cells and again in the
// draw pass to hand them their boxes, so it must draw the same cells in the same
// order both runs.
func (c *CSS) TableRow(draw func(*CSS)) {
	b := c.table
	if b == nil {
		draw(c)
		return
	}
	if c.measure {
		r := &tableRow{}
		c.rowWithin(b, r, draw)
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
	c.rowWithin(b, child.row, draw)
}

// rowWithin runs the cells of one row: the callback draws them and the row
// collects them. A row that is not the one in place is put back afterwards, so
// a table inside a cell cannot leave the outer row collecting.
func (c *CSS) rowWithin(b *tableBatch, r *tableRow, draw func(*CSS)) {
	outer := c.row
	c.row = r
	r.next = 0
	draw(c)
	c.row = outer
}

// layoutTable routes a widget drawn inside a Table callback. A widget drawn
// inside a row is a cell of it: the silent pass measures the cell and the draw
// pass hands back the box the solver gave it. A widget drawn straight in the
// table is a row of its own, the way an anonymous row works in a browser.
func (c *CSS) layoutTable(role, label string) (x, y, w, h int, st css.Style, ok bool) {
	if c.row == nil {
		return c.tableRowOfOne(role, label)
	}
	return c.layoutCell(role, label)
}

// tableRowOfOne wraps a widget drawn straight in a table in the row it belongs
// to: a row of a single cell, which is what a browser's anonymous boxes make of
// a child that did not ask to be a row. A widget that asked to be a caption is
// not wrapped — it is painted above the table, across its whole width.
func (c *CSS) tableRowOfOne(role, label string) (x, y, w, h int, st css.Style, ok bool) {
	b := c.table
	if c.measure {
		style := c.style.baseStyle(role, State{})
		if style.Caption() {
			cell := c.tableCellMeasure("", "", style)
			b.children = append(b.children, tableChild{caption: cell})
			return 0, 0, 0, 0, style, false
		}
		r := &tableRow{}
		c.rowWithin(b, r, func(*CSS) { c.layoutCell(role, label) })
		b.children = append(b.children, tableChild{row: r})
		return 0, 0, 0, 0, style, false
	}
	if b.rowIdx >= len(b.children) {
		return 0, 0, 0, 0, css.Style{}, false
	}
	child := b.children[b.rowIdx]
	b.rowIdx++
	if child.caption != nil {
		return c.placeCaption(child.caption)
	}
	return c.rowWithinDraw(child.row, role, label)
}

// rowWithinDraw hands one row its cells back in the draw pass.
func (c *CSS) rowWithinDraw(r *tableRow, role, label string) (x, y, w, h int, st css.Style, ok bool) {
	if r == nil {
		return 0, 0, 0, 0, st, false
	}
	c.rowWithin(c.table, r, func(*CSS) {
		x, y, w, h, st, ok = c.layoutCell(role, label)
	})
	return x, y, w, h, st, ok
}

// layoutCell places one cell of the row in place. The silent pass measures it
// and reports ok=false so nothing paints; the draw pass hands back the box the
// solver gave it.
func (c *CSS) layoutCell(role, label string) (x, y, w, h int, st css.Style, ok bool) {
	r := c.row
	b := c.table
	if c.measure {
		cell := c.tableCellMeasure(role, label, c.style.baseStyle(role, State{}))
		r.cells = append(r.cells, cell)
		return 0, 0, 0, 0, cell.st, false
	}
	e := c.style.beginWidget(role, label)
	st = c.style.entryStyle(e, role, e.state)
	if st.Display == css.DisplayNone {
		return 0, 0, 0, 0, st, false
	}
	if r.next >= len(r.cells) {
		return 0, 0, 0, 0, st, false
	}
	cell := r.cells[r.next]
	r.next++
	x, y = b.originX+cell.x, b.originY+cell.y
	return x, y, cell.w, cell.h, st, true
}

// tableCellMeasure reads one cell's style and measures the box it asks for:
// the natural size of the widget, capped by a width or height the cell's own
// style declares and clamped by its min and max. A percentage answers to the
// widest the table can be, which is the room the flow has for it — the table's
// own width is decided afterwards, by the columns.
func (c *CSS) tableCellMeasure(role, label string, st css.Style) *tableCell {
	cell := &tableCell{role: role, label: label, st: st}
	cell.nw, cell.nh = c.style.natural(role, label, st, c.style.u())
	winW := c.win.Width()
	if st.Has("width") {
		cell.nw = c.style.length(st, st.Width, winW)
	}
	if st.Has("max-width") {
		cell.nw = min(cell.nw, c.style.length(st, st.MaxWidth, winW))
	}
	if st.Has("min-width") {
		cell.nw = max(cell.nw, c.style.length(st, st.MinWidth, winW))
	}
	winH := c.win.Height()
	if st.Has("height") {
		cell.nh = c.style.length(st, st.Height, winH)
	}
	if st.Has("max-height") {
		cell.nh = min(cell.nh, c.style.length(st, st.MaxHeight, winH))
	}
	if st.Has("min-height") {
		cell.nh = max(cell.nh, c.style.length(st, st.MinHeight, winH))
	}
	return cell
}

// placeCaption gives a caption the table's whole width and the run above the
// rows, with the height it asks for: a caption is not in a column, so it is
// never given one.
func (c *CSS) placeCaption(cell *tableCell) (x, y, w, h int, st css.Style, ok bool) {
	b := c.table
	e := c.style.beginWidget(cell.role, cell.label)
	st = c.style.entryStyle(e, cell.role, e.state)
	if st.Display == css.DisplayNone {
		return 0, 0, 0, 0, st, false
	}
	return b.originX, b.originY - cell.nh, max(cell.nw, b.contentW), cell.nh, st, true
}
