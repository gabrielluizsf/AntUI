package template

import (
	"sort"

	"github.com/gabrielluizsf/antui/template/css"
)

// flexItem is one child of a flex container as the solver sees it: the box
// the stylesheet wants, collected in the silent pass and turned into a final
// x,y,w,h in the draw one. Every field the layout mutates is reset fresh by
// solve from nw/nh, so the two passes always answer the same geometry.
type flexItem struct {
	role, label  string
	st           css.Style
	nw, nh       int // measured border-box size, untouched by the solver
	w, h         int // solved border-box size
	ml, mr       int // margins on the main axis
	mt, mb       int // margins on the cross axis
	order        int
	grow, shrink float64
	basis        css.Length
	alignSelf    uint8 // AlignAuto means follow the container
	explicitW    bool  // width declared: stretch must not widen it
	explicitH    bool  // height declared: stretch must not heighten it
	minMain      int   // explicit main-axis minimum, -1 when absent
	maxMain      int   // explicit main-axis maximum, -1 when absent
	base         int   // main-axis size before grow or shrink
	out          bool  // out of flow: positioned by insets, not by the solver
	inLeft       int   // out-of-flow insets, -1 when absent
	inTop        int
	inRight      int
	inBottom     int
	x, y         int // solved position relative to the content origin
}

// flexBatch is one flex layout: the container's style and box on the page,
// and the items collected while its children drew. The silent pass fills
// items; solve turns them into boxes; the draw pass hands those boxes back.
// Nested [CSS.Flex] blocks flatten into the enclosing batch, so a container
// is always drawn by the Flex that started it.
type flexBatch struct {
	css      *CSS
	st       css.Style
	x, y, w  int // container border box
	h        int // solved after content has a size
	originX  int // content-box top-left in window coordinates
	originY  int
	contentW int // content-box width, fixed by the flow
	contentH int // content-box height, solved by the distribution
	items    []flexItem
	next     int // the item the draw pass is handing back
}

// col reports whether the container's main axis runs vertically.
func (b *flexBatch) col() bool {
	return b.st.FlexDirection == css.FlexDirectionColumn || b.st.FlexDirection == css.FlexDirectionColumnReverse
}

// Flex lays the widgets its callback draws out as flex items inside one flex
// container: a block in the flow that widens to the page, styled by the flex
// role in the stylesheet, and growing to hold its children along the
// flex-direction, wrapped and justified per the container's flex properties.
// Child widgets are placed by the solver, not the top-to-bottom flow, and
// keep all their interaction — a button inside a row still clicks.
//
// The callback runs twice per frame: first silently, to measure every child's
// natural box, then again to draw and interact at the boxes the algorithm
// assigned. It must draw the same widgets in the same order both runs — which
// a stable loop over a data slice does — and an event a child raises is seen
// on the second run, so read it before the callback returns:
//
//	clicked := false
//	app.Flex(func(f *CSS) {
//		for _, label := range items {
//			if f.Button(label).Is(event.Button, event.Click) {
//				clicked = true
//			}
//		}
//	})
//	if clicked { /* the click really happened */ }
//
// A Flex drawn inside another Flex's callback flattens into the outer
// container: the nested block's own face is dropped and its children join
// the outer items, which keeps the two-pass measurement honest.
func (c *CSS) Flex(draw func(*CSS)) {
	if c.flexing {
		draw(c)
		return
	}

	st := c.style.baseStyle(css.RoleFlex, State{})
	if st.Display == css.DisplayNone {
		return
	}
	u := c.style.u()
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

	bl := st.BorderWidth[3] * u
	br := st.BorderWidth[1] * u
	bt := st.BorderWidth[0] * u
	bb := st.BorderWidth[2] * u
	pl := c.style.length(st, st.Padding[3], winW)
	pr := c.style.length(st, st.Padding[1], winW)
	pt := c.style.length(st, st.Padding[0], winW)
	pb := c.style.length(st, st.Padding[2], winW)

	b := &flexBatch{
		css:      c,
		st:       st,
		x:        x,
		y:        y,
		w:        w,
		originX:  x + bl + pl,
		contentW: max(w-br-bl-pr-pl, 0),
	}
	c.flex = b
	c.flexing = true

	c.measure = true
	draw(c)
	c.measure = false

	// A container that declares its height bounds the content box up front,
	// so the solver has a definite main or cross size to distribute against.
	if st.Has("height") {
		hh := c.clampDim(st, c.style.length(st, st.Height, winH), "height", winH)
		b.contentH = max(hh-bt-bb-pt-pb, 0)
	}
	b.solve()

	h := b.contentH + bt + bb + pt + pb
	if st.Has("height") {
		h = c.clampDim(st, c.style.length(st, st.Height, winH), "height", winH)
	}
	h = max(h, c.clampDim(st, 0, "height", winH))
	b.h = h
	// The content box every child draws against is exactly what the border box
	// leaves inside, so an explicit height or a tall row agrees with the flow.
	b.contentH = max(h-bt-bb-pt-pb, 0)
	b.originY = y + bt + pt

	if c.style.boxOn(st) && visible(st) {
		c.style.paintBox(c.win, st, x, y, w, h)
	}

	b.next = 0
	draw(c)

	c.flexing = false
	c.flex = nil

	c.cursorY = y + h
	c.lastMB = mb
}

// clampDim applies the min/max constraint of one axis to a container
// dimension, against the window like every other percentage.
func (c *CSS) clampDim(st css.Style, v int, axis string, base int) int {
	switch axis {
	case "width":
		if st.Has("min-width") {
			v = max(v, c.style.length(st, st.MinWidth, base))
		}
		if st.Has("max-width") {
			v = min(v, c.style.length(st, st.MaxWidth, base))
		}
	case "height":
		if st.Has("min-height") {
			v = max(v, c.style.length(st, st.MinHeight, base))
		}
		if st.Has("max-height") {
			v = min(v, c.style.length(st, st.MaxHeight, base))
		}
	}
	return v
}

// layoutFlex routes a widget drawn inside a Flex callback. The silent pass
// measures it and returns ok=false so nothing paints; the draw pass hands
// back the box the solver placed it in. Out-of-flow children are parked by
// their insets against the container's content box and skip distribution.
func (c *CSS) layoutFlex(role, label string) (x, y, w, h int, st css.Style, ok bool) {
	b := c.flex
	if c.measure {
		it := c.flexItemMeasure(role, label)
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
		return 0, 0, 0, 0, st, false
	}
	it := b.items[b.next]
	b.next++
	if it.out {
		return c.flexPlaceOut(st, it)
	}
	return b.originX + it.x, b.originY + it.y, it.w, it.h, st, true
}

// flexPlaceOut positions an out-of-flow child by its inset edges against the
// container's content box, left/top when given and right/bottom from the far
// edge otherwise.
func (c *CSS) flexPlaceOut(st css.Style, it flexItem) (x, y, w, h int, os css.Style, ok bool) {
	b := c.flex
	x, y = b.originX, b.originY
	if it.inLeft >= 0 {
		x += it.inLeft
	} else if it.inRight >= 0 {
		x += b.contentW - it.w - it.inRight
	}
	if it.inTop >= 0 {
		y += it.inTop
	} else if it.inBottom >= 0 {
		y += b.contentH - it.h - it.inBottom
	}
	return x, y, it.w, it.h, st, true
}

// flexItemMeasure reads one flex item's style and measures its box: the
// natural size a widget wants, capped by its own width/height and clamped by
// min/max, with the margins, order and flex factors the style declares. The
// flex properties a container reads are here; the box the widget paints is
// the same one it would have taken in the flow. Properties the style never
// set keep their CSS initial values, which do not always match the struct's
// zero value: flex-shrink is 1, flex-basis and align-self are auto.
func (c *CSS) flexItemMeasure(role, label string) flexItem {
	b := c.flex
	st := c.style.baseStyle(role, State{})
	u := c.style.u()
	nw, nh := c.style.natural(role, label, st, u)
	if st.Has("width") {
		nw = c.style.length(st, st.Width, b.contentW)
	}
	if st.Has("min-width") {
		nw = max(nw, c.style.length(st, st.MinWidth, b.contentW))
	}
	if st.Has("max-width") {
		nw = min(nw, c.style.length(st, st.MaxWidth, b.contentW))
	}
	if st.Has("height") {
		nh = c.style.length(st, st.Height, b.contentW)
	}
	if st.Has("min-height") {
		nh = max(nh, c.style.length(st, st.MinHeight, b.contentW))
	}
	if st.Has("max-height") {
		nh = min(nh, c.style.length(st, st.MaxHeight, b.contentW))
	}
	shrink := 1.0
	if st.Has("flex-shrink") {
		shrink = st.FlexShrink
	}
	basis := css.Auto()
	if st.Has("flex-basis") {
		basis = st.FlexBasis
	}
	self := uint8(css.AlignAuto)
	if st.Has("align-self") {
		self = st.AlignSelf
	}
	it := flexItem{
		role: role, label: label, st: st,
		nw: nw, nh: nh, w: nw, h: nh,
		ml:        c.style.length(st, st.Margin[3], b.contentW),
		mr:        c.style.length(st, st.Margin[1], b.contentW),
		mt:        c.style.length(st, st.Margin[0], b.contentW),
		mb:        c.style.length(st, st.Margin[2], b.contentW),
		order:     st.Order,
		grow:      st.FlexGrow,
		shrink:    shrink,
		basis:     basis,
		alignSelf: self,
		explicitW: st.Has("width"),
		explicitH: st.Has("height"),
		minMain:   -1,
		maxMain:   -1,
		inLeft:    -1,
		inTop:     -1,
		inRight:   -1,
		inBottom:  -1,
	}
	if st.OutOfFlow() {
		it.out = true
		if st.Has("left") {
			it.inLeft = c.style.length(st, st.Left, b.contentW)
		}
		if st.Has("right") {
			it.inRight = c.style.length(st, st.Right, b.contentW)
		}
		if st.Has("top") {
			it.inTop = c.style.length(st, st.Top, b.contentW)
		}
		if st.Has("bottom") {
			it.inBottom = c.style.length(st, st.Bottom, b.contentW)
		}
		return it
	}
	if b.col() {
		if st.Has("min-height") {
			it.minMain = c.style.length(st, st.MinHeight, b.contentW)
		}
		if st.Has("max-height") {
			it.maxMain = c.style.length(st, st.MaxHeight, b.contentW)
		}
		return it
	}
	if st.Has("min-width") {
		it.minMain = c.style.length(st, st.MinWidth, b.contentW)
	}
	if st.Has("max-width") {
		it.maxMain = c.style.length(st, st.MaxWidth, b.contentW)
	}
	return it
}

// solve turns the measured items into boxes: decide the axes, split lines,
// distribute the free space and answer where every item sits and how tall the
// content grew. It runs once per pass; every run recomputes from the measured
// nw/nh, so the draw pass agrees with the silent one.
func (b *flexBatch) solve() {
	if len(b.items) == 0 {
		return
	}
	if b.col() {
		b.solveColumn()
		return
	}
	b.solveRow()
}

// baseMain picks an item's base main-axis size: its flex-basis when the basis
// is definite, its natural size otherwise. A percentage basis needs a
// definite axis to resolve against, or it falls back to the natural size.
func (b *flexBatch) baseMain(it *flexItem, availMain int, definite bool) int {
	natural := it.nw
	if b.col() {
		natural = it.nh
	}
	base := natural
	if !it.basis.Auto() && !it.basis.None() {
		if !(it.basis.IsPct() && !definite) {
			base = b.css.style.length(b.st, it.basis, availMain)
		}
	}
	if it.minMain >= 0 {
		base = max(base, it.minMain)
	}
	if it.maxMain >= 0 {
		base = min(base, it.maxMain)
	}
	return base
}

// sorted returns the item indices in flex order: order ascending, ties keep
// the draw order.
func (b *flexBatch) sorted() []int {
	idx := make([]int, len(b.items))
	for i := range idx {
		idx[i] = i
	}
	sort.SliceStable(idx, func(i, j int) bool {
		return b.items[idx[i]].order < b.items[idx[j]].order
	})
	return idx
}

// gapPx resolves one gap length. Percentages measure the container's content
// width; CSS ties row-gap to the content height, but the height only exists
// once solving finishes, so the width stands in for both.
func (b *flexBatch) gapPx(l css.Length) int {
	return b.css.style.length(b.st, l, b.contentW)
}

// split wraps the ordered items into lines along the main axis, packing until
// the next item would cross the edge. Items too big for a line own it alone.
func (b *flexBatch) split(order []int, avail, gap int, wrap bool) [][]int {
	lines := [][]int{}
	cur := []int{}
	used := 0
	for _, idx := range order {
		it := &b.items[idx]
		om := it.base + it.ml + it.mr
		if wrap && len(cur) > 0 && used+om+gap > avail {
			lines = append(lines, cur)
			cur = []int{}
			used = 0
		}
		if len(cur) > 0 {
			used += gap
		}
		used += om
		cur = append(cur, idx)
	}
	return append(lines, cur)
}

// myMain reads an item's solved main-axis size.
func (b *flexBatch) myMain(it *flexItem) int {
	if b.col() {
		return it.h
	}
	return it.w
}

// setMain stores the solved main-axis size; the cross stays where it was.
func (b *flexBatch) setMain(it *flexItem, size int) {
	if b.col() {
		it.h = size
		return
	}
	it.w = size
}

// distribute runs the flex-grow and flex-shrink algorithm over one line:
// grow shares the leftover space proportionally, shrink takes space back from
// the fattest bases, both floors clamped by an explicit minimum.
func (b *flexBatch) distribute(line []int, avail, gap int) {
	used, grow := 0, 0.0
	shrinkK := 0.0
	for _, idx := range line {
		it := &b.items[idx]
		used += it.base + it.ml + it.mr
		grow += it.grow
		shrinkK += float64(it.base) * it.shrink
	}
	free := avail - gap*(len(line)-1) - used
	switch {
	case free > 0 && grow > 0:
		for _, idx := range line {
			it := &b.items[idx]
			extra := int(float64(free) * (it.grow / grow))
			size := it.base + extra
			if it.maxMain >= 0 {
				size = min(size, it.maxMain)
			}
			b.setMain(it, size)
		}
	case free < 0 && shrinkK > 0:
		want := -free
		for _, idx := range line {
			it := &b.items[idx]
			cut := int(float64(want) * (float64(it.base) * it.shrink / shrinkK))
			b.setMain(it, max(it.minMain, it.base-cut))
		}
	default:
		for _, idx := range line {
			it := &b.items[idx]
			b.setMain(it, it.base)
		}
	}
}

// justify walks one line along the main axis, turning the leftover space into
// the gaps justify-content asks for, and mirrors the whole line when the
// direction runs backwards.
func (b *flexBatch) justify(line []int, avail, gap int) {
	used := 0
	for _, idx := range line {
		it := &b.items[idx]
		used += b.myMain(it) + it.ml + it.mr
	}
	left := max(avail-gap*(len(line)-1)-used, 0)
	lead, between := 0, gap
	reverse := b.st.FlexDirection == css.FlexDirectionRowReverse || b.st.FlexDirection == css.FlexDirectionColumnReverse
	switch b.st.JustifyContent {
	case css.JustifyFlexEnd:
		lead = left
	case css.JustifyCenter:
		lead = left / 2
	case css.JustifySpaceBetween:
		if len(line) > 1 {
			between = gap + left/(len(line)-1)
		}
	case css.JustifySpaceAround:
		between = gap + left/len(line)
		lead = left / (2 * len(line))
	case css.JustifySpaceEvenly:
		between = gap + left/(len(line)+1)
		lead = left / (len(line) + 1)
	}
	pos := lead
	for _, idx := range line {
		it := &b.items[idx]
		if b.col() {
			it.y = pos + it.mt
		} else {
			it.x = pos + it.ml
		}
		pos += it.ml + it.mr + b.myMain(it) + between
	}
	if !reverse {
		return
	}
	for _, idx := range line {
		it := &b.items[idx]
		if b.col() {
			it.y = avail - it.y - it.h
		} else {
			it.x = avail - it.x - it.w
		}
	}
}

// lineCross measures one line's natural cross size: the tallest (or widest)
// item plus its cross margins.
func (b *flexBatch) lineCross(line []int) int {
	cross := 0
	for _, idx := range line {
		it := &b.items[idx]
		size, margin := it.nh, it.mt+it.mb
		if b.col() {
			size, margin = it.nw, it.ml+it.mr
		}
		cross = max(cross, size+margin)
	}
	return cross
}

// alignCross plants every item across the cross axis of its line, offset by
// where the line itself starts: stretch grows an auto-sized item to fill, and
// the other keywords park it against a line edge or in the middle.
func (b *flexBatch) alignCross(line []int, cross, off int) {
	for _, idx := range line {
		it := &b.items[idx]
		a := b.st.AlignItems
		if it.alignSelf != css.AlignAuto {
			a = it.alignSelf
		}
		size := it.nh
		if b.col() {
			size = it.nw
		}
		if a == css.AlignStretch {
			if b.col() && !it.explicitW {
				size = max(cross-it.ml-it.mr, 0)
				it.w = size
			}
			if !b.col() && !it.explicitH {
				size = max(cross-it.mt-it.mb, 0)
				it.h = size
			}
		}
		switch a {
		case css.AlignCenter:
			if b.col() {
				it.x = off + (cross-size)/2
			} else {
				it.y = off + (cross-size)/2
			}
		case css.AlignFlexEnd:
			if b.col() {
				it.x = off + cross - size - it.mr
			} else {
				it.y = off + cross - size - it.mb
			}
		default: // stretch, flex-start, baseline
			if b.col() {
				it.x = off + it.ml
			} else {
				it.y = off + it.mt
			}
		}
	}
}

// stackLines answers where each line starts along the cross axis and how big
// the content grew: align-content spreads the leftover across the lines and
// stretch grows them to fill a definite cross size. The offsets are kept
// separate from the items so alignCross can add the within-line alignment.
func (b *flexBatch) stackLines(lines [][]int, crosses []int, gap, avail int, definite bool) ([]int, int) {
	if len(lines) == 0 {
		return nil, 0
	}
	total := 0
	for i, c := range crosses {
		if i > 0 {
			total += gap
		}
		total += c
	}
	if definite && b.st.AlignContent == css.ContentStretch {
		extra := max(avail-total, 0) / len(lines)
		for i := range crosses {
			crosses[i] += extra
		}
		total = 0
		for i, c := range crosses {
			if i > 0 {
				total += gap
			}
			total += c
		}
	}
	left := 0
	if definite {
		left = max(avail-total, 0)
	}
	lead, between := 0, gap
	switch b.st.AlignContent {
	case css.ContentFlexEnd:
		lead = left
	case css.ContentCenter:
		lead = left / 2
	case css.ContentSpaceBetween:
		if len(lines) > 1 {
			between = gap + left/(len(lines)-1)
		}
	case css.ContentSpaceAround:
		between = gap + left/len(lines)
		lead = left / (2 * len(lines))
	case css.ContentSpaceEvenly:
		between = gap + left/(len(lines)+1)
		lead = left / (len(lines) + 1)
	}
	offsets := make([]int, len(lines))
	pos := lead
	for li := range lines {
		offsets[li] = pos
		if li < len(lines)-1 {
			pos += crosses[li] + between
		}
	}
	if definite {
		return offsets, avail
	}
	return offsets, total
}

// solveRow lays the items along a horizontal main axis: wrap, grow, justify,
// then cross-align each line and stack the lines with align-content.
func (b *flexBatch) solveRow() {
	for i := range b.items {
		it := &b.items[i]
		it.base = b.baseMain(it, b.contentW, true)
	}
	order := b.sorted()
	wrap := b.st.FlexWrap != css.FlexWrapNowrap
	mainGap := b.gapPx(b.st.ColumnGap)
	lines := b.split(order, b.contentW, mainGap, wrap)

	for _, line := range lines {
		b.distribute(line, b.contentW, mainGap)
		b.justify(line, b.contentW, mainGap)
	}

	crosses := make([]int, len(lines))
	for li, line := range lines {
		crosses[li] = b.lineCross(line)
	}

	definite := b.st.Has("height")
	avail := 0
	if definite {
		avail = b.contentH
	}
	offsets, content := b.stackLines(lines, crosses, b.gapPx(b.st.RowGap), avail, definite)
	b.contentH = content

	for li, line := range lines {
		b.alignCross(line, crosses[li], offsets[li])
	}
}

// solveColumn lays the items along a vertical main axis, stacking wrapped
// groups next to each other. Without a declared height the axis is unbounded:
// nothing grows, shrinks or leans on justify-content, the column is a plain
// vertical stack of its natural heights.
func (b *flexBatch) solveColumn() {
	definite := b.st.Has("height")
	avail := b.contentH
	for i := range b.items {
		it := &b.items[i]
		it.base = b.baseMain(it, avail, definite)
	}
	order := b.sorted()
	wrap := definite && b.st.FlexWrap != css.FlexWrapNowrap
	mainGap := b.gapPx(b.st.RowGap)

	total := 0
	prev := -1
	for _, idx := range order {
		it := &b.items[idx]
		if prev >= 0 {
			total += mainGap
		}
		total += it.base + it.mt + it.mb
		prev = idx
	}

	eff := avail
	if !definite || eff < total {
		eff = total
	}
	lines := b.split(order, eff, mainGap, wrap)

	for _, line := range lines {
		b.distribute(line, eff, mainGap)
		b.justify(line, eff, mainGap)
	}

	crosses := make([]int, len(lines))
	for li, line := range lines {
		crosses[li] = b.lineCross(line)
	}
	offsets, _ := b.stackLines(lines, crosses, b.gapPx(b.st.ColumnGap), b.contentW, true)

	for li, line := range lines {
		b.alignCross(line, crosses[li], offsets[li])
	}

	if definite {
		b.contentH = avail
	} else {
		b.contentH = total
	}
}
