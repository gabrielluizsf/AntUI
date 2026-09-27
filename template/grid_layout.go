package template

import (
	"sort"

	"github.com/gabrielluizsf/antui/template/css"
)

type gridItem struct {
	role, label string
	st          css.Style
	nw, nh      int
	contribW    int
	contribH    int
	contribMinW int
	contribMinH int
	w, h        int
	ml, mr      int
	mt, mb      int
	explicitW   bool
	explicitH   bool
	minW        int
	maxW        int
	minH        int
	maxH        int
	column      css.GridPlacement
	row         css.GridPlacement
	justifySelf uint8
	alignSelf   uint8
	order       int
	colStart    int
	colEnd      int
	rowStart    int
	rowEnd      int
	colDefinite bool
	rowDefinite bool
	out         bool
	inLeft      int
	inTop       int
	inRight     int
	inBottom    int
	x, y        int
	nested      *gridNested // set when the item is a grid container of its own
}

type gridBatch struct {
	css      *CSS
	st       css.Style
	x, y     int
	w, h     int
	originX  int
	originY  int
	contentW int
	contentH int
	rowGap   int
	colGap   int
	items    []gridItem
	next     int
	// The explicit column tracks the last solve placed the items on, kept so a
	// container that has to size itself can ask how wide its own columns want
	// to be without resolving the template a second time.
	trackColumns []css.GridTrack
}

type gridArea struct {
	col      int
	row      int
	colSpan  int
	rowSpan  int
	hasOther bool
	otherCol int
	otherRow int
}

// Grid lays the widgets its callback draws out as grid items inside one grid
// container: a block in the flow that widens to the page, styled by the grid
// role in the stylesheet, and split into tracks per the container's
// grid-template properties. Child widgets are placed by the solver, not the
// top-to-bottom flow, and keep all their interaction — a button in a cell still
// clicks.
//
// The callback runs twice per frame: first silently, to measure every child's
// natural box, then again to draw and interact at the boxes the algorithm
// assigned. It must draw the same widgets in the same order both runs — which
// a stable loop over a data slice does — and an event a child raises is seen on
// the second run, so read it before the callback returns:
//
//	selected := 0
//	app.Grid(func(f *CSS) {
//		for i, cell := range cells {
//			if f.Button(cell.label).Is(event.Button, event.Click) {
//				selected = i
//			}
//		}
//	})
//
// Tracks are the fixed lengths, auto, fr and minmax() the grammar allows, plus
// repeat() — with a count, or with auto-fill and auto-fit to work out the count
// from the room the container has — and the content keywords min-content,
// max-content and fit-content(). Lines may be named in a template and placed
// by name, in either direction, and grid-template-areas gives a name to a
// rectangle of cells an item can ask for by grid-area.
//
// A grid drawn inside another grid's callback, or inside a flex callback, is an
// item of that container: the outer solver picks the cell it lands in, the
// container keeps its own face, and its children are laid out inside the box it
// was given, so a percentage or an fr track in it answers to the cell it ended
// up in. Containers carry no classes, so a nested grid is styled by the same
// grid rule as the one around it — a width declared there is a width on both.
//
//	subgrid is not supported: a nested grid lays its own tracks out from
//	scratch, and asking it to adopt the parent's lines is an error.
func (c *CSS) Grid(draw func(*CSS)) {
	if c.griding {
		c.nestedGrid(draw)
		return
	}
	if c.flexing {
		c.flexNestedGrid(draw)
		return
	}

	st := c.style.baseStyle(css.RoleGrid, State{})
	if st.Display == css.DisplayNone {
		return
	}
	winW, winH := c.win.Width(), c.win.Height()
	ml := c.style.length(st, st.Margin[3], winW)
	mr := c.style.length(st, st.Margin[1], winW)
	mt := c.style.length(st, st.Margin[0], winW)
	mb := c.style.length(st, st.Margin[2], winW)

	c.cursorY += c.lineH
	c.lineH = 0
	c.lineX = 0
	c.cursorY += max(c.lastMB, mt)
	y := c.cursorY

	w := max(winW-ml-mr, 0)
	if st.Has("width") && gridLengthDefinite(st.Width) {
		w = c.style.length(st, st.Width, winW)
	}
	w = c.clampDim(st, w, "width", winW)
	x := ml + (winW-ml-mr-w)/2

	b := &gridBatch{css: c, st: st}
	c.setGridBox(b, st, x, y, w, c.gridDeclaredHeight(st), winW)
	c.grid, c.griding = b, true
	c.measure = true
	draw(c)
	c.measure = false
	b.solve()

	_, _, top, bottom := c.gridChrome(st, winW)
	if st.Has("min-height") && gridLengthDefinite(st.MinHeight) {
		limit := c.clampDim(st, c.style.length(st, st.MinHeight, winH), "height", winH)
		b.contentH = max(b.contentH, max(limit-top-bottom, 0))
	}
	if st.Has("max-height") && gridLengthDefinite(st.MaxHeight) {
		limit := c.clampDim(st, c.style.length(st, st.MaxHeight, winH), "height", winH)
		b.contentH = min(b.contentH, max(limit-top-bottom, 0))
	}
	h := b.contentH + top + bottom
	if st.Has("height") && gridLengthDefinite(st.Height) {
		h = c.style.length(st, st.Height, winH)
	}
	h = c.clampDim(st, h, "height", winH)
	c.setGridBox(b, st, x, y, w, h, winW)

	if c.style.boxOn(st) && visible(st) {
		c.style.paintBox(c.win, st, x, y, w, h)
	}

	c.gridWithin(b, draw)

	c.grid = nil
	c.griding = false
	c.cursorY = y + h
	c.lastMB = mb
}

// gridChrome is the room a container's own border and padding take out of its
// border box, with a percentage padding answered to pctBase.
func (c *CSS) gridChrome(st css.Style, pctBase int) (left, right, top, bottom int) {
	u := c.style.u()
	return st.BorderWidth[3]*u + c.style.length(st, st.Padding[3], pctBase),
		st.BorderWidth[1]*u + c.style.length(st, st.Padding[1], pctBase),
		st.BorderWidth[0]*u + c.style.length(st, st.Padding[0], pctBase),
		st.BorderWidth[2]*u + c.style.length(st, st.Padding[2], pctBase)
}

// setGridBox gives a batch the geometry of the container's border box on
// screen: where it sits, how big it is, and the content box its children are
// placed in. pctBase is what a percentage inside the box answers to, and a
// percentage row-gap is answered to what a declared height leaves, since a
// container sized by its content has no height to divide.
func (c *CSS) setGridBox(b *gridBatch, st css.Style, x, y, w, h, pctBase int) {
	left, right, top, bottom := c.gridChrome(st, pctBase)
	b.x, b.y, b.w, b.h = x, y, w, h
	b.originX, b.originY = x+left, y+top
	b.contentW = max(w-left-right, 0)
	b.contentH = max(h-top-bottom, 0)
	rowBase := b.contentH
	if st.Has("height") && gridLengthDefinite(st.Height) {
		rowBase = max(c.style.length(st, st.Height, c.win.Height())-top-bottom, 0)
	}
	b.rowGap = c.style.length(st, st.RowGap, rowBase)
	b.colGap = c.style.length(st, st.ColumnGap, b.contentW)
}

// gridDeclaredHeight is the border-box height the stylesheet gives a container,
// clamped by its own min and max, and 0 when it declares none. It bounds the
// content box before the rows are sized, which is what gives them a space to
// fill and a definite axis to divide.
func (c *CSS) gridDeclaredHeight(st css.Style) int {
	if !st.Has("height") || !gridLengthDefinite(st.Height) {
		return 0
	}
	winH := c.win.Height()
	return max(c.clampDim(st, c.style.length(st, st.Height, winH), "height", winH), 0)
}

// gridWithin runs the children of one grid container: the callback draws them,
// the silent pass collects them and the draw pass hands them their boxes. The
// batch in place is the container they belong to, so a grid nested in this one
// is collected as its item, and a container of another kind has its own solver
// put aside while these children are placed.
func (c *CSS) gridWithin(b *gridBatch, draw func(*CSS)) {
	outer, outerGriding, outerFlexing := c.grid, c.griding, c.flexing
	c.grid, c.griding, c.flexing = b, true, false
	b.next = 0
	draw(c)
	c.grid, c.griding, c.flexing = outer, outerGriding, outerFlexing
}

// layoutGrid routes a widget drawn inside a Grid callback. The silent pass
// measures it and returns ok=false so nothing paints; the draw pass hands back
// the box the solver placed it in. An item the flow left out is parked by its
// placement — the area grid-area names, or the padding box — and its insets.
func (c *CSS) layoutGrid(role, label string) (x, y, w, h int, st css.Style, ok bool) {
	b := c.grid
	if c.measure {
		it := c.gridItemMeasure(role, label)
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
		return c.gridPlaceOut(st, it)
	}
	x = b.originX + it.x
	y = b.originY + it.y
	x, y = c.shiftInset(st, x, y, b.contentW, b.contentH)
	return x, y, it.w, it.h, st, true
}

func (c *CSS) gridItemMeasure(role, label string) gridItem {
	st := c.style.baseStyle(role, State{})
	u := c.style.u()
	naturalW, naturalH := c.style.natural(role, label, st, u)
	// A track asking for min-content wants the narrowest the item can be, and
	// the size the item declares or constrains is the narrowest it will ever be.
	minNaturalW, minNaturalH := c.style.minContent(role, label, st, u)
	return c.gridItemSized(role, label, st, naturalW, naturalH, minNaturalW, minNaturalH)
}

// gridItemSized collects one item the way the stylesheet asks for it, from the
// size the widget wants and the narrowest it can be. The solver gets to see
// both, so a track can ask for either.
func (c *CSS) gridItemSized(role, label string, st css.Style, naturalW, naturalH, minNaturalW, minNaturalH int) gridItem {
	b := c.grid
	nw, nh := naturalW, naturalH
	contribW, contribH := naturalW, naturalH
	minContribW, minContribH := minNaturalW, minNaturalH
	rowDefinite := b.st.Has("height") && gridLengthDefinite(b.st.Height)
	if st.Has("width") && gridLengthDefinite(st.Width) {
		nw = c.style.length(st, st.Width, b.contentW)
		if !st.Width.IsPct() {
			contribW, minContribW = nw, nw
		}
	}
	if st.Has("max-width") && gridLengthDefinite(st.MaxWidth) {
		limit := c.style.length(st, st.MaxWidth, b.contentW)
		nw = min(nw, limit)
		if !st.MaxWidth.IsPct() {
			contribW = min(contribW, limit)
			minContribW = min(minContribW, limit)
		}
	}
	if st.Has("min-width") && gridLengthDefinite(st.MinWidth) {
		limit := c.style.length(st, st.MinWidth, b.contentW)
		nw = max(nw, limit)
		if !st.MinWidth.IsPct() {
			contribW = max(contribW, limit)
			minContribW = max(minContribW, limit)
		}
	}
	if st.Has("height") && gridLengthDefinite(st.Height) && (!st.Height.IsPct() || rowDefinite) {
		nh = c.style.length(st, st.Height, b.contentH)
		if !st.Height.IsPct() {
			contribH, minContribH = nh, nh
		}
	}
	if st.Has("max-height") && gridLengthDefinite(st.MaxHeight) && (!st.MaxHeight.IsPct() || rowDefinite) {
		limit := c.style.length(st, st.MaxHeight, b.contentH)
		nh = min(nh, limit)
		if !st.MaxHeight.IsPct() {
			contribH = min(contribH, limit)
			minContribH = min(minContribH, limit)
		}
	}
	if st.Has("min-height") && gridLengthDefinite(st.MinHeight) && (!st.MinHeight.IsPct() || rowDefinite) {
		limit := c.style.length(st, st.MinHeight, b.contentH)
		nh = max(nh, limit)
		if !st.MinHeight.IsPct() {
			contribH = max(contribH, limit)
			minContribH = max(minContribH, limit)
		}
	}

	justifySelf := uint8(css.AlignAuto)
	if st.Has("justify-self") {
		justifySelf = st.JustifySelf
	}
	alignSelf := uint8(css.AlignAuto)
	if st.Has("align-self") {
		alignSelf = st.AlignSelf
	}
	it := gridItem{
		role:        role,
		label:       label,
		st:          st,
		nw:          nw,
		nh:          nh,
		contribW:    contribW,
		contribH:    contribH,
		contribMinW: minContribW,
		contribMinH: minContribH,
		w:           nw,
		h:           nh,
		ml:          c.style.length(st, st.Margin[3], b.contentW),
		mr:          c.style.length(st, st.Margin[1], b.contentW),
		mt:          c.style.length(st, st.Margin[0], b.contentW),
		mb:          c.style.length(st, st.Margin[2], b.contentW),
		explicitW:   st.Has("width") && gridLengthDefinite(st.Width),
		explicitH:   st.Has("height") && gridLengthDefinite(st.Height),
		minW:        -1,
		maxW:        -1,
		minH:        -1,
		maxH:        -1,
		column:      st.GridColumn,
		row:         st.GridRow,
		justifySelf: justifySelf,
		alignSelf:   alignSelf,
		order:       st.Order,
		inLeft:      -1,
		inTop:       -1,
		inRight:     -1,
		inBottom:    -1,
	}
	if st.Has("min-width") && gridLengthDefinite(st.MinWidth) {
		it.minW = c.style.length(st, st.MinWidth, b.contentW)
	}
	if st.Has("max-width") && gridLengthDefinite(st.MaxWidth) {
		it.maxW = c.style.length(st, st.MaxWidth, b.contentW)
	}
	if st.Has("min-height") && gridLengthDefinite(st.MinHeight) && (!st.MinHeight.IsPct() || rowDefinite) {
		it.minH = c.style.length(st, st.MinHeight, b.contentH)
	}
	if st.Has("max-height") && gridLengthDefinite(st.MaxHeight) && (!st.MaxHeight.IsPct() || rowDefinite) {
		it.maxH = c.style.length(st, st.MaxHeight, b.contentH)
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
	}
	return it
}

// gridPlaceOut hands an out-of-flow item the box solve worked out for it.
func (c *CSS) gridPlaceOut(st css.Style, it gridItem) (int, int, int, int, css.Style, bool) {
	return c.grid.originX + it.x, c.grid.originY + it.y, it.w, it.h, st, true
}

// solve lays every item out: it resolves both axes, places each item on them,
// then hands the axes their space and turns the placements into boxes.
func (b *gridBatch) solve() {
	// A container that declares its height bounds the content box up front, so
	// the row axis has a definite size to divide and a definite space for the
	// rows an auto-repeat counts into.
	rowDefinite := b.st.Has("height") && gridLengthDefinite(b.st.Height)
	rowBase := b.contentW
	if rowDefinite {
		rowBase = b.contentH
	}
	columnAxis := b.resolveGridAxis(b.st.GridTemplateColumns, b.contentW, b.colGap)
	rowAxis := b.resolveGridAxis(b.st.GridTemplateRows, b.gridRowSpace(rowDefinite), b.rowGap)
	columns := columnAxis.Tracks
	rows := rowAxis.Tracks
	areaColumns, areaRows := 0, 0
	for _, row := range b.st.GridTemplateAreas {
		areaColumns = max(areaColumns, len(row))
	}
	areaRows = len(b.st.GridTemplateAreas)
	ensureGridTracks(&columns, areaColumns)
	ensureGridTracks(&rows, areaRows)
	columnShift := b.gridLineShift(columnAxis, false)
	rowShift := b.gridLineShift(rowAxis, true)
	prependGridTracks(&columns, columnShift)
	prependGridTracks(&rows, rowShift)
	areas := b.gridAreas()
	occupied := make(map[[2]int]bool, len(b.items)*2)

	for i := range b.items {
		it := &b.items[i]
		it.w, it.h = it.nw, it.nh
		if area, ok := areas[it.gridAreaName()]; ok {
			it.colStart, it.colEnd = area.col+columnShift, area.col+area.colSpan+columnShift
			it.rowStart, it.rowEnd = area.row+rowShift, area.row+area.rowSpan+rowShift
			it.colDefinite, it.rowDefinite = true, true
			ensureGridTracks(&columns, it.colEnd)
			ensureGridTracks(&rows, it.rowEnd)
			if !it.out {
				b.occupy(occupied, it)
			}
			continue
		}
		it.colStart, it.colEnd, it.colDefinite = gridAxisPlacement(columnAxis, it.column)
		it.colStart += columnShift
		it.colEnd += columnShift
		it.rowStart, it.rowEnd, it.rowDefinite = gridAxisPlacement(rowAxis, it.row)
		it.rowStart += rowShift
		it.rowEnd += rowShift
		if it.out || !(it.colDefinite && it.rowDefinite) {
			continue
		}
		ensureGridTracks(&columns, it.colEnd)
		ensureGridTracks(&rows, it.rowEnd)
		b.occupy(occupied, it)
	}

	order := b.gridOrder()
	rowCursor, colCursor := 0, 0
	place := func(index int) {
		it := &b.items[index]
		if !it.out && !(it.colDefinite && it.rowDefinite) {
			b.placeAutomatic(it, &columns, &rows, occupied, &rowCursor, &colCursor)
		}
	}
	// The items locked to a line come first, in the order they were added, and
	// everything else follows in that same order: an item with a definite column
	// and an automatic one shares the pass, so it takes the cursor where the
	// items before it left it instead of jumping ahead of them.
	if b.st.GridAutoFlow == css.GridAutoFlowRow {
		for _, index := range order {
			it := &b.items[index]
			if !it.out && it.rowDefinite && !it.colDefinite {
				place(index)
			}
		}
		for _, index := range order {
			it := &b.items[index]
			if !it.out && !(it.rowDefinite && !it.colDefinite) {
				place(index)
			}
		}
	} else {
		for _, index := range order {
			it := &b.items[index]
			if !it.out && it.colDefinite && !it.rowDefinite {
				place(index)
			}
		}
		for _, index := range order {
			it := &b.items[index]
			if !it.out && !(it.colDefinite && !it.rowDefinite) {
				place(index)
			}
		}
	}

	rowGapTotal := max(len(rows)-1, 0) * b.rowGap
	colGapTotal := max(len(columns)-1, 0) * b.colGap
	columnSizes := b.resolveGridTracks(columns, true, true, max(b.contentW-colGapTotal, 0))
	rowSizes := b.resolveGridTracks(rows, false, rowDefinite, max(rowBase-rowGapTotal, 0))
	if !rowDefinite {
		b.contentH = sumGridSizes(rowSizes) + rowGapTotal
	}
	b.trackColumns = columns

	columnPositions := gridTrackPositions(columnSizes, b.colGap, sumGridSizes(columnSizes), b.contentW, b.st.GridJustifyContent)
	rowPositions := gridTrackPositions(rowSizes, b.rowGap, sumGridSizes(rowSizes), b.contentH, b.st.AlignContent)
	for i := range b.items {
		if b.items[i].out {
			b.placeOutOfFlow(&b.items[i], columnSizes, rowSizes, columnPositions, rowPositions)
			continue
		}
		b.placeGridItem(&b.items[i], columnSizes, rowSizes, columnPositions, rowPositions)
	}
}

// placeOutOfFlow positions an item the flow left out, in coordinates from the
// content origin like every other item. The box it is measured against is the
// area its placement names, and the padding box when the placement names none;
// the insets then move it inside that box, and a pair of opposite insets with
// no width of its own fills what is left between them.
func (b *gridBatch) placeOutOfFlow(it *gridItem, columns, rows, columnPositions, rowPositions []int) {
	baseX, baseW := 0, b.contentW
	if it.colDefinite {
		baseX, baseW = gridAreaBox(columns, columnPositions, it.colStart, it.colEnd, b.colGap)
	}
	baseY, baseH := 0, b.contentH
	if it.rowDefinite {
		baseY, baseH = gridAreaBox(rows, rowPositions, it.rowStart, it.rowEnd, b.rowGap)
	}

	it.w, it.h = it.nw, it.nh
	// A pair of opposite insets on an axis whose size was not declared fills
	// what is left between them, so the width is settled before the position.
	if !it.explicitW && it.inLeft >= 0 && it.inRight >= 0 {
		it.w = max(baseW-it.inLeft-it.inRight, 0)
	}
	if !it.explicitH && it.inTop >= 0 && it.inBottom >= 0 {
		it.h = max(baseH-it.inTop-it.inBottom, 0)
	}
	it.x, it.y = baseX, baseY
	if it.inLeft >= 0 {
		it.x = baseX + it.inLeft
	} else if it.inRight >= 0 {
		it.x = baseX + baseW - it.inRight - it.w
	}
	if it.inTop >= 0 {
		it.y = baseY + it.inTop
	} else if it.inBottom >= 0 {
		it.y = baseY + baseH - it.inBottom - it.h
	}
	if it.minW >= 0 {
		it.w = max(it.w, it.minW)
	}
	if it.maxW >= 0 {
		it.w = min(it.w, it.maxW)
	}
	if it.minH >= 0 {
		it.h = max(it.h, it.minH)
	}
	if it.maxH >= 0 {
		it.h = min(it.h, it.maxH)
	}
}

// gridAreaBox is the stretch of the tracks an item spans: where it starts on
// screen and how wide it is, gaps included.
func gridAreaBox(sizes, positions []int, start, end, gap int) (int, int) {
	if len(sizes) == 0 {
		return 0, 0
	}
	start = min(max(start, 0), len(sizes)-1)
	end = min(max(end, start+1), len(sizes))
	width := max(end-start-1, 0) * gap
	for track := start; track < end; track++ {
		width += sizes[track]
	}
	return positions[start], width
}

// resolveGridAxis turns a parsed template into the explicit tracks of one axis,
// working out how many repetitions an auto-fill or auto-fit stands for and
// carrying the line names along with them. The rows of a container that does
// not declare its height have no space to fill, so an auto-repeat there keeps
// the single track the grammar promises.
func (b *gridBatch) resolveGridAxis(template css.GridTemplate, available, gap int) css.GridAxis {
	// A track minimum is a length in the sheet's own pixels, so the space it is
	// counted against has to be in those too. A window drawn at twice its size
	// measures twice as wide on the canvas, and without the conversion an
	// auto-repeat would fit twice as many tracks as the room holds.
	unit := max(Scale(b.css.win), 1)
	return template.Resolve(available/unit, gap/unit, b.gridInFlow())
}

// gridRowSpace is the space a row axis can fill: the declared content height
// when there is one, and nothing at all when the axis is sized by its content.
func (b *gridBatch) gridRowSpace(definite bool) int {
	if !definite {
		return 0
	}
	return b.contentH
}

// gridInFlow is how many items the auto-placement has to find room for, which
// is what an auto-fit counts its repetitions against.
func (b *gridBatch) gridInFlow() int {
	count := 0
	for i := range b.items {
		if !b.items[i].out {
			count++
		}
	}
	return count
}

func (b *gridBatch) gridAreas() map[string]gridArea {
	areas := make(map[string]gridArea)
	for row, cells := range b.st.GridTemplateAreas {
		for col, name := range cells {
			if name == "." || name == "..." {
				continue
			}
			area, exists := areas[name]
			if !exists {
				areas[name] = gridArea{col: col, row: row, colSpan: 1, rowSpan: 1}
				continue
			}
			area.colSpan = max(area.col, col) - min(area.col, col) + 1
			area.rowSpan = max(area.row, row) - min(area.row, row) + 1
			if area.col != col || area.row != row {
				area.hasOther = true
				area.otherCol, area.otherRow = col, row
			}
			areas[name] = area
		}
	}
	return areas
}

// gridAreaName is the name a grid-area gave all four of its lines at once,
// which is how an item finds the template area it was pointed at.
func (it *gridItem) gridAreaName() string {
	if it.column.Start.Kind != css.GridLineName || it.column.End.Kind != css.GridLineName || it.row.Start.Kind != css.GridLineName || it.row.End.Kind != css.GridLineName {
		return ""
	}
	name := it.column.Start.Name
	if name == "" || it.column.End.Name != name || it.row.Start.Name != name || it.row.End.Name != name {
		return ""
	}
	return name
}

// gridAxisPlacement reads a placement off one axis into the range of tracks it
// covers. A line is an index counted from the start or the end, or a name the
// template carries — the first line carrying the name for a start, the last for
// an end, which is how a name bounds both sides of an item. A name no line
// carries is no placement at all, leaving the item to be placed automatically.
func gridAxisPlacement(axis css.GridAxis, placement css.GridPlacement) (int, int, bool) {
	start, startOK := gridAxisLine(axis, placement.Start, false)
	end, endOK := gridAxisLine(axis, placement.End, true)
	if !startOK && !endOK {
		return 0, 0, false
	}

	if !startOK {
		start = end - gridPlacementSpan(placement)
	}
	if !endOK {
		if placement.End.Kind == css.GridLineSpan {
			end = start + placement.End.Index
		} else if placement.Start.Span > 0 {
			end = start + placement.Start.Span
		} else {
			end = start + 1
		}
	}
	if end < start {
		start, end = end, start
	}
	if end <= start {
		end = start + 1
	}
	return start, end, true
}

// gridAxisLine is the line one side of a placement names, and whether the
// placement named one at all. A count is taken from the start unless it is
// negative, and a name is looked up from the end of the axis when the end of
// the item is what named it or when a "-" sent it there.
func gridAxisLine(axis css.GridAxis, line css.GridLine, end bool) (int, bool) {
	explicit := len(axis.Tracks)
	switch line.Kind {
	case css.GridLineIndex:
		if line.Index > 0 {
			return line.Index - 1, true
		}
		if line.Index < 0 {
			return explicit + line.Index + 1, true
		}
	case css.GridLineName:
		if found, ok := axis.LineIndex(line.Name, end || line.Backward); ok {
			return found, true
		}
	}
	return 0, false
}

func gridPlacementSpan(placement css.GridPlacement) int {
	if placement.Start.Kind == css.GridLineSpan {
		return max(1, placement.Start.Index)
	}
	if placement.End.Kind == css.GridLineSpan {
		return max(1, placement.End.Index)
	}
	return 1
}

func ensureGridTracks(tracks *[]css.GridTrack, count int) {
	for len(*tracks) < count {
		*tracks = append(*tracks, css.GridTrack{Kind: css.GridTrackSingle, Size: css.GridTrackSize{Kind: css.GridTrackAuto}})
	}
}

func prependGridTracks(tracks *[]css.GridTrack, count int) {
	if count <= 0 {
		return
	}
	prefix := make([]css.GridTrack, count)
	for i := range prefix {
		prefix[i] = css.GridTrack{Kind: css.GridTrackSingle, Size: css.GridTrackSize{Kind: css.GridTrackAuto}}
	}
	*tracks = append(prefix, *tracks...)
}

// gridLineShift is how many implicit lines have to stand before the explicit
// ones, for the items that count lines from the end of the axis.
func (b *gridBatch) gridLineShift(axis css.GridAxis, rows bool) int {
	shift := 0
	for i := range b.items {
		placement := b.items[i].column
		if rows {
			placement = b.items[i].row
		}
		start, end, ok := gridAxisPlacement(axis, placement)
		if ok {
			shift = min(shift, start, end)
		}
	}
	return max(-shift, 0)
}

func (b *gridBatch) occupy(occupied map[[2]int]bool, it *gridItem) {
	for row := it.rowStart; row < it.rowEnd; row++ {
		for col := it.colStart; col < it.colEnd; col++ {
			occupied[[2]int{col, row}] = true
		}
	}
}

func (b *gridBatch) fits(occupied map[[2]int]bool, col, row, colSpan, rowSpan int) bool {
	if col < 0 || row < 0 || colSpan < 1 || rowSpan < 1 {
		return false
	}
	for r := row; r < row+rowSpan; r++ {
		for c := col; c < col+colSpan; c++ {
			if occupied[[2]int{c, r}] {
				return false
			}
		}
	}
	return true
}

func (b *gridBatch) reserve(occupied map[[2]int]bool, it *gridItem, col, row, colSpan, rowSpan int) {
	it.colStart, it.colEnd = col, col+colSpan
	it.rowStart, it.rowEnd = row, row+rowSpan
	it.colDefinite, it.rowDefinite = true, true
	b.occupy(occupied, it)
}

func (b *gridBatch) gridOrder() []int {
	order := make([]int, len(b.items))
	for i := range order {
		order[i] = i
	}
	sort.SliceStable(order, func(i, j int) bool {
		return b.items[order[i]].order < b.items[order[j]].order
	})
	return order
}

func (b *gridBatch) placeAutomatic(it *gridItem, columns *[]css.GridTrack, rows *[]css.GridTrack, occupied map[[2]int]bool, rowCursor, colCursor *int) {
	colSpan, rowSpan := gridPlacementSpan(it.column), gridPlacementSpan(it.row)
	if it.colDefinite {
		colSpan = max(it.colEnd-it.colStart, colSpan)
	}
	if it.rowDefinite {
		rowSpan = max(it.rowEnd-it.rowStart, rowSpan)
	}

	if b.st.GridAutoFlow == css.GridAutoFlowColumn {
		b.placeColumnAutomatic(it, columns, rows, occupied, rowCursor, colCursor, colSpan, rowSpan)
		return
	}
	b.placeRowAutomatic(it, columns, rows, occupied, rowCursor, colCursor, colSpan, rowSpan)
}

func (b *gridBatch) placeRowAutomatic(it *gridItem, columns *[]css.GridTrack, rows *[]css.GridTrack, occupied map[[2]int]bool, rowCursor, colCursor *int, colSpan, rowSpan int) {
	if it.colDefinite {
		startRow := 0
		if !b.st.GridAutoFlowDense {
			startRow = *rowCursor
			if it.colStart < *colCursor {
				startRow++
			}
		}
		for row := startRow; ; row++ {
			ensureGridTracks(rows, row+rowSpan)
			if b.fits(occupied, it.colStart, row, colSpan, rowSpan) {
				b.reserve(occupied, it, it.colStart, row, colSpan, rowSpan)
				*rowCursor, *colCursor = row, it.colStart+colSpan
				return
			}
		}
	}
	if it.rowDefinite {
		startCol := 0
		if !b.st.GridAutoFlowDense && it.rowStart == *rowCursor {
			startCol = *colCursor
		}
		for col := startCol; ; col++ {
			ensureGridTracks(columns, col+colSpan)
			if b.fits(occupied, col, it.rowStart, colSpan, rowSpan) {
				b.reserve(occupied, it, col, it.rowStart, colSpan, rowSpan)
				*rowCursor, *colCursor = it.rowStart, col+colSpan
				return
			}
		}
	}
	if b.st.GridAutoFlowDense {
		for row := 0; row < len(*rows); row++ {
			for col := 0; col < len(*columns); col++ {
				if b.fits(occupied, col, row, colSpan, rowSpan) {
					b.reserve(occupied, it, col, row, colSpan, rowSpan)
					return
				}
			}
		}
	}
	row, col := *rowCursor, *colCursor
	for {
		if col > 0 && col+colSpan > len(*columns) {
			col = 0
			row++
			*rowCursor, *colCursor = row, col
			continue
		}
		ensureGridTracks(columns, col+colSpan)
		ensureGridTracks(rows, row+rowSpan)
		if b.fits(occupied, col, row, colSpan, rowSpan) {
			b.reserve(occupied, it, col, row, colSpan, rowSpan)
			*rowCursor, *colCursor = row, col+colSpan
			return
		}
		col++
	}
}

func (b *gridBatch) placeColumnAutomatic(it *gridItem, columns *[]css.GridTrack, rows *[]css.GridTrack, occupied map[[2]int]bool, rowCursor, colCursor *int, colSpan, rowSpan int) {
	if it.rowDefinite {
		startCol := 0
		if !b.st.GridAutoFlowDense && it.rowStart == *rowCursor {
			startCol = *colCursor
		}
		for col := startCol; ; col++ {
			ensureGridTracks(columns, col+colSpan)
			if b.fits(occupied, col, it.rowStart, colSpan, rowSpan) {
				b.reserve(occupied, it, col, it.rowStart, colSpan, rowSpan)
				*rowCursor, *colCursor = it.rowStart, col+colSpan
				return
			}
		}
	}
	if it.colDefinite {
		startRow := 0
		if !b.st.GridAutoFlowDense {
			startRow = *rowCursor
			if it.colStart < *colCursor {
				startRow++
			}
		}
		for row := startRow; ; row++ {
			ensureGridTracks(rows, row+rowSpan)
			if b.fits(occupied, it.colStart, row, colSpan, rowSpan) {
				b.reserve(occupied, it, it.colStart, row, colSpan, rowSpan)
				*rowCursor, *colCursor = row+rowSpan, it.colStart
				return
			}
		}
	}
	if b.st.GridAutoFlowDense {
		for col := 0; col < len(*columns); col++ {
			for row := 0; row < len(*rows); row++ {
				if b.fits(occupied, col, row, colSpan, rowSpan) {
					b.reserve(occupied, it, col, row, colSpan, rowSpan)
					return
				}
			}
		}
	}
	col, row := *colCursor, *rowCursor
	for {
		if row > 0 && row+rowSpan > len(*rows) {
			row = 0
			col++
			*rowCursor, *colCursor = row, col
			continue
		}
		ensureGridTracks(columns, col+colSpan)
		ensureGridTracks(rows, row+rowSpan)
		if b.fits(occupied, col, row, colSpan, rowSpan) {
			b.reserve(occupied, it, col, row, colSpan, rowSpan)
			*rowCursor, *colCursor = row+rowSpan, col
			return
		}
		row++
	}
}

func (b *gridBatch) placeGridItem(it *gridItem, columns, rows, columnPositions, rowPositions []int) {
	if it.colStart < 0 || it.colStart >= len(columns) || it.rowStart < 0 || it.rowStart >= len(rows) {
		return
	}
	cellX := columnPositions[it.colStart]
	cellY := rowPositions[it.rowStart]
	cellW, cellH := 0, 0
	for col := it.colStart; col < it.colEnd && col < len(columns); col++ {
		cellW += columns[col]
	}
	for row := it.rowStart; row < it.rowEnd && row < len(rows); row++ {
		cellH += rows[row]
	}
	cellW += max(it.colEnd-it.colStart-1, 0) * b.colGap
	cellH += max(it.rowEnd-it.rowStart-1, 0) * b.rowGap

	justify := b.st.JustifyItems
	if !b.st.Has("justify-items") {
		justify = css.AlignStretch
	}
	if it.justifySelf != css.AlignAuto {
		justify = it.justifySelf
	}
	align := b.st.AlignItems
	if !b.st.Has("align-items") {
		align = css.AlignStretch
	}
	if it.alignSelf != css.AlignAuto {
		align = it.alignSelf
	}

	it.w, it.h = it.nw, it.nh
	if it.st.Width.IsPct() {
		it.w = b.css.style.length(it.st, it.st.Width, cellW)
	}
	if it.st.Height.IsPct() {
		it.h = b.css.style.length(it.st, it.st.Height, cellH)
	}
	if justify == css.AlignStretch && !it.explicitW {
		it.w = max(cellW-it.ml-it.mr, 0)
	}
	if align == css.AlignStretch && !it.explicitH {
		it.h = max(cellH-it.mt-it.mb, 0)
	}
	minW, maxW := it.minW, it.maxW
	if it.st.MinWidth.IsPct() {
		minW = b.css.style.length(it.st, it.st.MinWidth, cellW)
	}
	if it.st.MaxWidth.IsPct() {
		maxW = b.css.style.length(it.st, it.st.MaxWidth, cellW)
	}
	minH, maxH := it.minH, it.maxH
	if it.st.MinHeight.IsPct() {
		minH = b.css.style.length(it.st, it.st.MinHeight, cellH)
	}
	if it.st.MaxHeight.IsPct() {
		maxH = b.css.style.length(it.st, it.st.MaxHeight, cellH)
	}
	it.w = clampGridItem(it.w, minW, maxW)
	it.h = clampGridItem(it.h, minH, maxH)

	switch justify {
	case css.AlignCenter:
		it.x = cellX + (cellW-it.ml-it.mr-it.w)/2 + it.ml
	case css.AlignFlexEnd:
		it.x = cellX + cellW - it.mr - it.w
	default:
		it.x = cellX + it.ml
	}
	switch align {
	case css.AlignCenter:
		it.y = cellY + (cellH-it.mt-it.mb-it.h)/2 + it.mt
	case css.AlignFlexEnd:
		it.y = cellY + cellH - it.mb - it.h
	default:
		it.y = cellY + it.mt
	}
}

func gridLengthDefinite(length css.Length) bool {
	return !length.Auto() && !length.None()
}

func clampGridItem(value, minimum, maximum int) int {
	if maximum >= 0 {
		value = min(value, maximum)
	}
	if minimum >= 0 {
		value = max(value, minimum)
	}
	return value
}

func gridTrackPositions(sizes []int, gap, total, available int, alignment uint8) []int {
	positions := make([]int, len(sizes))
	if len(sizes) == 0 {
		return positions
	}
	used := total + gap*(len(sizes)-1)
	left := max(available-used, 0)
	lead, between := 0, gap
	switch alignment {
	case css.ContentFlexEnd:
		lead = left
	case css.ContentCenter:
		lead = left / 2
	case css.ContentSpaceBetween:
		if len(sizes) > 1 {
			between += left / (len(sizes) - 1)
		}
	case css.ContentSpaceAround:
		between += left / len(sizes)
		lead = left / (2 * len(sizes))
	case css.ContentSpaceEvenly:
		between += left / (len(sizes) + 1)
		lead = left / (len(sizes) + 1)
	}
	position := lead
	for i := range sizes {
		positions[i] = position
		position += sizes[i] + between
	}
	return positions
}
