package template

import "github.com/gabrielluizsf/antui/template/css"

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
