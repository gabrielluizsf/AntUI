package template

import "github.com/gabrielluizsf/antui/template/css"

// A table's geometry is three questions asked in order: how wide is each
// column, how wide is the whole table, and how tall is each row. The columns
// ask the cells — a column is as wide as the widest cell that lands in it — the
// table asks the columns, and a table with a width of its own shares that width
// out over them. The rows ask the cells too, and a cell's content is placed in
// the row by the vertical-align its own style declares.

// tableChrome is the border and padding the table's own style adds around the
// cells, which a percentage inside it answers to.
func (c *CSS) tableChrome(st css.Style, pctBase int) (left, right, top, bottom int) {
	u := c.style.u()
	return st.BorderWidth[3]*u + c.style.length(st, st.Padding[3], pctBase),
		st.BorderWidth[1]*u + c.style.length(st, st.Padding[1], pctBase),
		st.BorderWidth[0]*u + c.style.length(st, st.Padding[0], pctBase),
		st.BorderWidth[2]*u + c.style.length(st, st.Padding[2], pctBase)
}

// tableColumns is the natural width of each column: the widest cell that lands
// in it, over every row of the table.
func (c *CSS) tableColumns(b *tableBatch) []int {
	var cols []int
	for _, ch := range b.children {
		if ch.row == nil {
			continue
		}
		for j, cell := range ch.row.cells {
			for len(cols) <= j {
				cols = append(cols, 0)
			}
			cols[j] = max(cols[j], cell.nw)
		}
	}
	return cols
}

// tableWidth is how wide the table's content box is: the width the stylesheet
// gave it, or the width its columns ask for when it gave none — which is what a
// table does when nothing told it how wide to be. The width never exceeds the
// room the flow has for it.
func (c *CSS) tableWidth(b *tableBatch, cols []int, avail, chrome int) int {
	w := colsWidth(cols)
	if b.st.Has("width") {
		w = max(c.clampDim(b.st, c.style.length(b.st, b.st.Width, avail), "width", avail), 0)
	}
	return min(w, max(avail-chrome, 0))
}

// tableShare is the width each column ends up with. The whole width is shared
// out in proportion to what the columns asked for, so a wide cell keeps more of
// the extra room than a narrow one, and the last column takes the rounding
// remainder so the columns add up to the width exactly.
func tableShare(cols []int, width int) []int {
	if len(cols) == 0 {
		return nil
	}
	natural := 0
	for _, w := range cols {
		natural += w
	}
	out := make([]int, len(cols))
	if natural == 0 {
		base, extra := width/len(cols), width%len(cols)
		for i := range out {
			out[i] = base
			if i < extra {
				out[i]++
			}
		}
		return out
	}
	shared := 0
	for i, w := range cols {
		if i == len(cols)-1 {
			out[i] = max(width-shared, 0)
			break
		}
		out[i] = width * w / natural
		shared += out[i]
	}
	return out
}

// solveTable gives every cell the width of its column and every row the height
// of its tallest cell, then gives the table the box that makes: a caption above
// the rows, the rows themselves and the table's own border and padding.
//
// x0 and room are the box the table is placed in — the page for a table in the
// flow, the cell, column or flex item for a nested one — and the table centres
// in it the way a block box centres in the flow. pctBase is what a percentage of
// the table itself answers to.
func (c *CSS) solveTable(b *tableBatch, x0, room, pctBase int) {
	winH := c.win.Height()
	left, right, top, bottom := c.tableChrome(b.st, pctBase)
	ml := c.style.length(b.st, b.st.Margin[3], pctBase)
	mr := c.style.length(b.st, b.st.Margin[1], pctBase)
	avail := max(room-ml-mr, 0)

	caption := 0
	for _, ch := range b.children {
		if ch.caption != nil {
			caption = max(ch.caption.nh, 0)
		}
	}
	b.caption = caption

	// The columns ask the cells what they need, and the table's own width — the
	// one it declared, or the one they add up to — is shared back out over them.
	natural := c.tableColumns(b)
	shared := tableShare(natural, c.tableWidth(b, natural, avail, left+right))
	b.cols = shared

	rows := 0
	for _, ch := range b.children {
		if ch.row == nil {
			continue
		}
		height := 0
		for _, cell := range ch.row.cells {
			height = max(height, cell.nh)
		}
		ch.row.height = height
		rows += height
	}

	h := caption + rows + top + bottom
	if b.st.Has("height") && gridLengthDefinite(b.st.Height) {
		h = c.clampDim(b.st, c.style.length(b.st, b.st.Height, winH), "height", winH)
	}
	h = max(h, c.clampDim(b.st, 0, "height", winH))

	w := colsWidth(shared) + left + right
	c.setTableBox(b, b.st, x0+ml+(avail-w)/2, b.y, w, h, pctBase)

	rowY := 0
	for _, ch := range b.children {
		if ch.row == nil {
			continue
		}
		x := 0
		for j, cell := range ch.row.cells {
			cell.w = 0
			if j < len(shared) {
				cell.w = shared[j]
			}
			cell.x = x
			cell.h, cell.y = cellPlacement(cell, ch.row.height, rowY)
			x += cell.w
		}
		rowY += ch.row.height
	}
}

// cellPlacement is the box a cell takes in its row: the whole row's height for
// a cell that asks for the top, and its own height placed in the middle or at
// the bottom for one that asks for that.
func cellPlacement(cell *tableCell, rowHeight, rowY int) (h, y int) {
	switch cell.st.VerticalAlign {
	case css.VerticalAlignMiddle, css.VerticalAlignSub, css.VerticalAlignSuper:
		return cell.nh, rowY + (rowHeight-cell.nh)/2
	case css.VerticalAlignBottom, css.VerticalAlignTextBottom:
		return cell.nh, rowY + max(rowHeight-cell.nh, 0)
	default:
		return rowHeight, rowY
	}
}

// colsWidth is what the columns add up to.
func colsWidth(cols []int) int {
	total := 0
	for _, w := range cols {
		total += w
	}
	return total
}
