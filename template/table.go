package template

import "github.com/gabrielluizsf/antui/template/css"

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
	mt := c.style.length(st, st.Margin[css.SideTop], winW)
	mb := c.style.length(st, st.Margin[css.SideBottom], winW)

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
