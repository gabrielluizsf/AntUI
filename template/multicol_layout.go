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

// multiPiece is one fragment: the column it lands in, how far down that column
// it starts and how tall the run of it is.
type multiPiece struct {
	col int
	y   int
	h   int
}

// multiPieces bounds the fragments one box may be cut into. Each fragment is
// one pass over the container's callback, and a box that would need more than
// this — a very tall box in a very short column — is drawn whole in the piece
// that overflows instead, so a stylesheet cannot cost a frame a thousand
// passes.
const multiPieces = 64

// multiBatch is one multi-column container: its box on the page, the columns
// the flow made, and the items collected while its children drew. The silent
// pass measures every child, the flow cuts the boxes into pieces, and the draw
// passes hand those pieces back — one pass per fragment an item was cut into,
// so every widget is drawn once per piece it owns.
type multiBatch struct {
	css      *CSS
	st       css.Style
	x, y, w  int // container border box
	h        int // solved after the columns have a height
	originX  int // content-box top-left in window coordinates
	originY  int
	contentW int
	contentH int
	left     int // the border and padding beside the columns
	right    int
	gap      int  // the gutter between two columns
	count    int  // how many columns there are
	colW     int  // the width of each one
	top      int  // the border and padding above the columns
	bottom   int  // the border and padding below them
	room     int  // the height one column may hold
	colH     int  // the tallest column's content
	cuts     bool // a height decided outside the flow: a box that does not fit is cut
	definite bool // the height of the box was decided, not asked for
	items    []multiItem
	next     int // the item the draw pass is handing back
	pass     int // which fragment pass the draw is on
	passes   int // how many passes the pieces need
}

// MultiCol lays the widgets its callback draws out in columns: a block in the
// flow that widens to the page, styled by the multicol role, whose content is
// split into equal columns with a gutter between them and an optional rule
// drawn in the gutter. The number of columns is column-count, the width that
// fits is column-width, the gutter is column-gap and the rule is column-rule —
// the same gutter flex and grid read as their gap, so one declaration spaces
// all three.
//
//	app.MultiCol(func(m *CSS) {
//		for _, item := range items {
//			m.Label(item)
//		}
//	})
//
// A container with no height grows to the tallest of its columns and balances
// the content over them, so column-count: 2 with five widgets puts two in the
// first column and three in the second. A container that declares a height
// spends it: the flow cuts a box that runs out of room and continues it in the
// next column, and a box with break-inside: avoid is moved whole to the next
// column instead of being cut. Content that does not fit the columns the
// container asked for spills into extra columns to the right, the way a browser
// spills a paragraph past a fixed column-count.
//
// The callback runs once silently to measure every child, then once per
// fragment a child was cut into, so it must draw the same widgets in the same
// order every run. Events are read on the last run, as [CSS.Flex] does.
func (c *CSS) MultiCol(draw func(*CSS)) {
	if c.multicoling || c.flexing || c.griding || c.tabling {
		c.nestedMultiCol(draw)
		return
	}

	st := c.style.baseStyle(css.RoleMultiCol, State{})
	if st.Display == css.DisplayNone {
		return
	}
	winW, winH := c.win.Width(), c.win.Height()
	ml := c.style.length(st, st.Margin[3], winW)
	mr := c.style.length(st, st.Margin[1], winW)
	mt := c.style.length(st, st.Margin[0], winW)
	mb := c.style.length(st, st.Margin[2], winW)

	// The container is a block: it closes the inline line, collapses the
	// sibling margins and rests on the flow cursor.
	c.cursorY += c.lineH
	c.lineH = 0
	c.lineX = 0
	c.cursorY += max(c.lastMB, mt)
	y := c.cursorY

	w := max(winW-ml-mr, 0)
	if st.Has("width") {
		w = c.style.length(st, st.Width, winW)
	}
	w = c.clampDim(st, w, "width", winW)
	x := ml + (winW-ml-mr-w)/2

	b := &multiBatch{css: c, st: st}
	c.setMultiBox(b, st, x, y, w, c.multiDeclaredHeight(st), winW)
	c.multi, c.multicoling = b, true

	c.measure = true
	draw(c)
	c.measure = false
	b.flow()

	h := max(b.colH+b.top+b.bottom, c.clampDim(st, 0, "height", winH))
	if b.cuts {
		h = c.multiDeclaredHeight(st)
	}
	b.h = h
	b.contentH = max(h-b.top-b.bottom, 0)

	if c.style.boxOn(st) && visible(st) {
		c.style.paintBox(c.win, st, x, y, w, h)
	}
	c.style.paintColumnRules(c.win, b)

	for pass := range b.passes {
		b.pass = pass
		c.multiWithin(b, draw)
	}

	c.multicoling = false
	c.multi = nil
	c.cursorY = y + h
	c.lastMB = mb
}

// setMultiBox gives a batch the geometry of the container's border box: where
// it sits, how wide it is, and the content box the columns are cut out of.
// The columns are solved here, before the children are measured, so a width in
// a percentage answers to the column it will be drawn in. pctBase is what a
// percentage of the container itself answers to, and a height of zero means the
// container is as tall as its columns turn out to be.
func (c *CSS) setMultiBox(b *multiBatch, st css.Style, x, y, w, h, pctBase int) {
	u := c.style.u()
	left := st.BorderWidth[3]*u + c.style.length(st, st.Padding[3], pctBase)
	right := st.BorderWidth[1]*u + c.style.length(st, st.Padding[1], pctBase)
	b.top = st.BorderWidth[0]*u + c.style.length(st, st.Padding[0], pctBase)
	b.bottom = st.BorderWidth[2]*u + c.style.length(st, st.Padding[2], pctBase)
	b.x, b.y, b.w, b.h = x, y, w, h
	b.originX, b.originY = x+left, y+b.top
	b.left, b.right = left, right
	b.contentW = max(w-left-right, 0)
	b.contentH = max(h-b.top-b.bottom, 0)
	b.gap = c.style.length(st, st.ColumnGap, b.contentW)
	b.solveColumns(st)
	// A height a stylesheet or a parent decided is one the flow can cut boxes
	// against. A height of zero is the container asking for whatever its columns
	// come to, which is the height nestedMultiCol's measure pass is about to
	// produce, and a parent that then hands that height back must not be read as
	// a decision — drawMultiContainer is the one that settles that.
	b.definite = h > 0
}

// multiWithin runs the children of one multi-column container: the callback
// draws them, the silent pass collects them and the draw passes hand them their
// pieces. The batch in place is the container they belong to, so another
// container inside this one is collected as its item, and a container of
// another kind has its own solver put aside while these children are placed.
func (c *CSS) multiWithin(b *multiBatch, draw func(*CSS)) {
	outer, outerMulti, outerFlex, outerGrid, outerTable := c.multi, c.multicoling, c.flexing, c.griding, c.tabling
	c.multi, c.multicoling, c.flexing, c.griding, c.tabling = b, true, false, false, false
	b.next = 0
	draw(c)
	c.multi, c.multicoling, c.flexing, c.griding, c.tabling = outer, outerMulti, outerFlex, outerGrid, outerTable
}

// layoutMulti routes a widget drawn inside a MultiCol callback. The silent pass
// measures it and returns ok=false so nothing paints; a draw pass hands back
// the piece the flow cut for the pass it is on.
func (c *CSS) layoutMulti(role, label string) (x, y, w, h int, st css.Style, ok bool) {
	b := c.multi
	if c.measure {
		it := c.multiItemMeasure(role, label)
		if it.st.Display == css.DisplayNone {
			return 0, 0, 0, 0, it.st, false
		}
		b.items = append(b.items, it)
		return 0, 0, 0, 0, it.st, false
	}

	e := c.style.beginWidget(role, label)
	st = c.style.entryStyle(e, role, e.state)
	if st.Display == css.DisplayNone {
		return 0, 0, 0, 0, st, false
	}
	if b.next >= len(b.items) {
		return 0, 0, 0, 0, st, true
	}
	it := &b.items[b.next]
	b.next++
	piece := b.pieceAt(it)
	// A child of a column is a block in it: it takes the column's width less
	// its own margins, and a piece starts inside its top margin.
	x = b.originX + piece.col*(b.colW+b.gap) + it.ml
	y = b.originY + piece.y
	w = max(b.colW-it.ml-it.mr, 0)
	return x, y, w, piece.h, st, true
}

// pieceAt is the fragment the draw pass is on for one item: the piece of the
// pass's own number, and the item's last one once the flow has nothing further
// to cut.
func (b *multiBatch) pieceAt(it *multiItem) multiPiece {
	p := min(b.pass, len(it.pieces)-1)
	if p < 0 {
		p = 0
	}
	return it.pieces[p]
}

// multiItemMeasure reads one child's style and measures its box: the natural
// size a widget wants, with the margins the style declares. The width and
// height it ends up with are resolved by the flow, which knows the column it
// will be drawn in.
func (c *CSS) multiItemMeasure(role, label string) multiItem {
	st := c.style.baseStyle(role, State{})
	nw, nh := c.style.natural(role, label, st, c.style.u())
	return c.multiItemSized(st, nw, nh)
}

// multiItemSized collects one child the way the stylesheet asks for it, from a
// size the child wants: a widget's own natural box, or the box a nested
// container asked for. The properties a multi-column flow reads are the same
// either way, so a container drawn inside this one is placed as a real child.
func (c *CSS) multiItemSized(st css.Style, nw, nh int) multiItem {
	b := c.multi
	return multiItem{
		st: st,
		nw: nw, nh: nh, w: nw, h: nh,
		ml:   c.style.length(st, st.Margin[3], b.colW),
		mr:   c.style.length(st, st.Margin[1], b.colW),
		mt:   c.style.length(st, st.Margin[0], b.colW),
		mb:   c.style.length(st, st.Margin[2], b.colW),
		keep: st.BreakInside == css.BreakAvoid,
	}
}

// multiItemBox is the width and height a child takes in a column colW wide: the
// size it asked for, capped by max and lifted by min. The flow calls it for
// every solve, so a width in a percentage answers to the column the child
// really ended up in.
func (c *CSS) multiItemBox(st css.Style, nw, nh, colW int) (w, h int) {
	if st.Has("width") {
		nw = c.style.length(st, st.Width, colW)
	}
	if st.Has("max-width") {
		nw = min(nw, c.style.length(st, st.MaxWidth, colW))
	}
	if st.Has("min-width") {
		nw = max(nw, c.style.length(st, st.MinWidth, colW))
	}
	if st.Has("height") {
		nh = c.style.length(st, st.Height, colW)
	}
	if st.Has("max-height") {
		nh = min(nh, c.style.length(st, st.MaxHeight, colW))
	}
	if st.Has("min-height") {
		nh = max(nh, c.style.length(st, st.MinHeight, colW))
	}
	return nw, nh
}
