package template

// A grid container drawn inside another container is an item of it, not a
// flattening of its children into the parent: the parent's solver picks the
// cell the container lands in, the container keeps its own face, and its
// children are laid out inside the box it was given.
//
// The silent pass measures the container where the parent would stretch it —
// the width a percentage or a flexible track inside it answers to — and then
// re-solves it at the width it asks for, which is what the parent is told. The
// draw pass re-solves it once more, at the box the parent really gave it, so a
// percentage or an fr track answers to the cell the container ended up in.

import "github.com/gabrielluizsf/antui/template/css"

// gridNested is one grid container drawn inside another container. It carries
// the batch that lays the container's children out — collected by the silent
// pass and re-solved at the final box — and the border box the container asks
// its parent for.
type gridNested struct {
	batch      *gridBatch
	nw, nh     int
	minW, minH int
}

// nestedGrid is [CSS.Grid] called from inside another grid: the container
// becomes an item of the parent, placed like any other and holding its own
// children.
func (c *CSS) nestedGrid(draw func(*CSS)) {
	st := c.style.baseStyle(css.RoleGrid, State{})
	if st.Display == css.DisplayNone {
		return
	}
	parent := c.grid
	if c.measure {
		nested := c.measureGridContainer(st, draw, parent.contentW)
		it := c.gridItemSized("", "", st, nested.nw, nested.nh, nested.minW, nested.minH)
		it.nested = nested
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
	c.drawGridContainer(it.nested, draw, parent.contentW, parent.originX+it.x, parent.originY+it.y, it.w, it.h)
}

// flexNestedGrid is [CSS.Grid] called from inside a flex: the container becomes
// a flex item, sized and placed by the flex solver like any other.
func (c *CSS) flexNestedGrid(draw func(*CSS)) {
	st := c.style.baseStyle(css.RoleGrid, State{})
	if st.Display == css.DisplayNone {
		return
	}
	parent := c.flex
	if c.measure {
		nested := c.measureGridContainer(st, draw, parent.contentW)
		it := c.flexItemSized(st, nested.nw, nested.nh)
		it.nested = nested
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
	c.drawGridContainer(it.nested, draw, parent.contentW, parent.originX+it.x, parent.originY+it.y, it.w, it.h)
}

// measureGridContainer lays a grid container out on its own, so its parent has
// a box to place: the width its own columns ask for, the height that width
// produces, and the narrowest of both. pctBase is what a percentage inside the
// container answers to, which is the width the parent would stretch it to.
func (c *CSS) measureGridContainer(st css.Style, draw func(*CSS), pctBase int) *gridNested {
	left, right, top, bottom := c.gridChrome(st, pctBase)
	height := c.gridDeclaredHeight(st)
	ml := c.style.length(st, st.Margin[3], pctBase)
	mr := c.style.length(st, st.Margin[1], pctBase)
	width := max(pctBase-ml-mr, 0)
	if st.Has("width") && gridLengthDefinite(st.Width) {
		width = max(c.style.length(st, st.Width, pctBase)-left-right, 0)
	}

	b := &gridBatch{css: c, st: st}
	c.setGridBox(b, st, 0, 0, width+left+right, height, pctBase)
	// The silent pass of the whole page is already running, so the children are
	// collected here and the container is sized from them.
	c.gridWithin(b, draw)
	b.solve()

	// The narrowest the container can be, and the width it asks for: its own
	// columns say how wide, and the rows answer to whatever width that is.
	minW, widest := b.gridIntrinsicWidth()
	c.setGridBox(b, st, 0, 0, minW+left+right, height, pctBase)
	b.solve()
	minH := b.contentH + top + bottom
	if widest != b.contentW {
		c.setGridBox(b, st, 0, 0, widest+left+right, height, pctBase)
		b.solve()
	}
	return &gridNested{
		batch: b,
		nw:    b.contentW + left + right,
		nh:    b.contentH + top + bottom,
		minW:  minW + left + right,
		minH:  minH,
	}
}

// drawGridContainer gives the container the box its parent solved for it,
// re-solves it at that size so the tracks answer to it, paints its face and
// draws its children inside.
func (c *CSS) drawGridContainer(nested *gridNested, draw func(*CSS), pctBase, x, y, w, h int) {
	b := nested.batch
	c.setGridBox(b, b.st, x, y, w, h, pctBase)
	b.solve()
	if c.style.boxOn(b.st) && visible(b.st) {
		c.style.paintBox(c.win, b.st, x, y, w, h)
	}
	c.gridWithin(b, draw)
}
