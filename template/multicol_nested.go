package template

import "github.com/gabrielluizsf/antui/template/css"

// A multi-column container drawn inside another container is an item of it, not
// a flattening of its children into the parent: the parent's solver picks the
// box the container lands in, the container keeps its own columns, and its
// children are laid out inside it.
//
// The silent pass measures the container where the parent would stretch it —
// the width a percentage or a column-width inside it answers to — and the draw
// pass re-solves it at the box the parent really gave it, so both answer to the
// column the container ended up in. A container in a column is never cut: the
// flow has no way to know where its own columns would be split, so it is laid
// out whole.

// multiNested is one multi-column container drawn inside another container. It
// carries the batch that lays the container's children out — collected by the
// silent pass and re-solved at the final box — and the border box the container
// asks its parent for.
type multiNested struct {
	batch  *multiBatch
	nw, nh int
}

// nestedMultiCol is [CSS.MultiCol] called from inside another container: the
// container becomes an item of the parent, placed like any other and holding
// its own children.
func (c *CSS) nestedMultiCol(draw func(*CSS)) {
	st := c.style.baseStyle(css.RoleMultiCol, State{})
	if st.Display == css.DisplayNone {
		return
	}
	switch {
	case c.multicoling:
		c.multiInMulti(st, draw)
	case c.flexing:
		c.flexNestedMultiCol(st, draw)
	case c.griding:
		c.gridNestedMultiCol(st, draw)
	default:
		// Straight in a table: the content of the cell in place, or of a row
		// of its own when nothing put one there.
		c.multiInTableCell(st, draw)
	}
}

func (c *CSS) multiInMulti(st css.Style, draw func(*CSS)) {
	parent := c.multi
	if c.measure {
		nested := c.measureMultiContainer(st, draw, parent.contentW)
		it := c.multiItemSized(st, nested.nw, nested.nh)
		it.keep, it.nested = true, nested
		parent.items = append(parent.items, it)
		return
	}
	if parent.next >= len(parent.items) {
		return
	}
	it := &parent.items[parent.next]
	if it.nested == nil {
		return
	}
	parent.next++
	piece := parent.pieceAt(it)
	x := parent.originX + it.ml
	y := parent.originY + piece.y
	c.drawMultiContainer(it.nested, draw, parent.colW, x, y, max(parent.colW-it.ml-it.mr, 0), it.h)
}

// tableInColumns puts a table inside one column of the multi-column container
// in place. The column is the room the table has, so the table shares its
// columns out over that width, and the flow never cuts it: there is no way to
// know where the table's own rows would be split, so it is laid out whole.
func (c *CSS) tableInColumns(st css.Style, draw func(*CSS)) {
	parent := c.multi
	if c.measure {
		nested := c.measureTableContainer(st, draw, parent.colW)
		it := c.multiItemSized(st, nested.nw, nested.nh)
		it.keep, it.table = true, nested
		parent.items = append(parent.items, it)
		return
	}
	if parent.next >= len(parent.items) {
		return
	}
	it := &parent.items[parent.next]
	parent.next++
	if it.table == nil {
		return
	}
	piece := parent.pieceAt(it)
	x := parent.originX + piece.col*(parent.colW+parent.gap) + it.ml
	y := parent.originY + piece.y
	room := max(parent.colW-it.ml-it.mr, 0)
	// A table is as wide as its columns ask for, not as wide as the column it
	// lands in, so it centres in the column the way any block box is centred.
	c.drawTableContainer(it.table, draw, room, x, y, it.h)
}

func (c *CSS) flexNestedMultiCol(st css.Style, draw func(*CSS)) {
	parent := c.flex
	if c.measure {
		nested := c.measureMultiContainer(st, draw, parent.contentW)
		it := c.flexItemSized(st, nested.nw, nested.nh)
		it.columns = nested
		parent.items = append(parent.items, it)
		return
	}
	if parent.next >= len(parent.items) {
		return
	}
	it := &parent.items[parent.next]
	if it.columns == nil {
		return
	}
	parent.next++
	c.drawMultiContainer(it.columns, draw, it.w, parent.originX+it.x, parent.originY+it.y, it.w, it.h)
}

func (c *CSS) gridNestedMultiCol(st css.Style, draw func(*CSS)) {
	parent := c.grid
	if c.measure {
		nested := c.measureMultiContainer(st, draw, parent.contentW)
		it := c.gridItemSized("", "", st, nested.nw, nested.nh, nested.nw, nested.nh)
		it.columns = nested
		parent.items = append(parent.items, it)
		return
	}
	if parent.next >= len(parent.items) {
		return
	}
	it := &parent.items[parent.next]
	if it.columns == nil {
		return
	}
	parent.next++
	c.drawMultiContainer(it.columns, draw, it.w, parent.originX+it.x, parent.originY+it.y, it.w, it.h)
}

// measureMultiContainer lays a multi-column container out on its own, so its
// parent has a box to place: the width the parent would stretch it to, and the
// height that width produces. pctBase is what a percentage of the container
// answers to. The children are collected once and re-flowed at the narrower
// width a parent that shrinks its items may end up giving.
func (c *CSS) measureMultiContainer(st css.Style, draw func(*CSS), pctBase int) *multiNested {
	ml := c.style.length(st, st.Margin[3], pctBase)
	mr := c.style.length(st, st.Margin[1], pctBase)
	width := max(pctBase-ml-mr, 0)
	if st.Has("width") {
		width = c.style.length(st, st.Width, pctBase)
	}
	width = c.clampDim(st, width, "width", pctBase)

	b := &multiBatch{css: c, st: st}
	c.setMultiBox(b, st, 0, 0, width, c.multiDeclaredHeight(st), pctBase)
	// The silent pass of the whole page is already running, so the children are
	// collected here and the container is sized from them.
	c.multiWithin(b, draw)
	b.flow()
	return &multiNested{batch: b, nw: b.ask(), nh: b.h}
}

// ask is how wide a column block being measured says it wants to be to whatever
// holds it. A block fills the room it is given, but the room a measure pass
// offers is a guess — in a cell whose columns are not sized it is the whole
// window — so a width in a percentage answers with what the columns inside need
// instead: one column's worth, or the widest box when there is a single column.
// That is what a percentage with no definite box to resolve against counts as,
// and it is what sizes the column that holds the block. The draw pass solves it
// again at the width the parent really gave.
func (b *multiBatch) ask() int {
	if !b.st.Has("width") || !b.st.Width.IsPct() {
		return b.w
	}
	w := b.colW
	if b.count < 2 {
		w = 0
		for i := range b.items {
			w = max(w, b.items[i].nw)
		}
	}
	return w + b.left + b.right
}

// drawMultiContainer gives the container the box its parent solved for it,
// re-solves it at that size so the columns answer to it, paints its face and
// its rules and draws its children inside.
func (c *CSS) drawMultiContainer(nested *multiNested, draw func(*CSS), pctBase, x, y, w, h int) {
	b := nested.batch
	// The box the parent really gave may be the one the container asked for, in
	// which case it keeps balancing over its columns. A different one — a row
	// taller than its content, a flex item stretched, a height of its own — is a
	// height to cut its boxes against.
	definite := c.multiDeclaredHeight(b.st) > 0 || h != nested.nh
	c.setMultiBox(b, b.st, x, y, w, h, pctBase)
	b.definite = definite
	b.flow()
	if c.style.boxOn(b.st) && visible(b.st) {
		c.style.paintBox(c.win, b.st, x, y, w, b.h)
	}
	c.style.paintColumnRules(c.win, b)
	for pass := range b.passes {
		b.pass = pass
		c.multiWithin(b, draw)
	}
}
