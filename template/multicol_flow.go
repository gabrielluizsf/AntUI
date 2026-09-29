package template

import "github.com/gabrielluizsf/antui/template/css"

// The flow of a multi-column container is one walk down the columns. Each
// child takes the room left in the column it is in; a child that runs out of
// room is cut and continues in the next column, unless it asked not to be cut,
// in which case it starts the next column whole. A container with no height
// balances instead: it measures what the children want, divides that by the
// column count and fills each column to the share, so the columns come out the
// same height and nothing is cut.

// multiDeclaredHeight is the border-box height the stylesheet gives a
// container, clamped by its own min and max, and 0 when it declares none.
func (c *CSS) multiDeclaredHeight(st css.Style) int {
	if !st.Has("height") || !gridLengthDefinite(st.Height) {
		return 0
	}
	winH := c.win.Height()
	return max(c.clampDim(st, c.style.length(st, st.Height, winH), "height", winH), 0)
}

// multiColumnRoom is how much content height one column may hold. A container
// that fills a column before starting the next spends its whole height on each
// one; a balancing container spreads the height over them, which is what a
// column-count with a declared height does.
func (c *CSS) multiColumnRoom(b *multiBatch, st css.Style) int {
	room := max(b.contentH, 0)
	if st.ColumnFill == css.ColumnFillAuto || b.count < 2 {
		return room
	}
	return max(room/b.count, 0)
}

// multiBalance is the height each column gets when the container has no height
// of its own: the tallest column of an even split of everything the children
// asked for. A child taller than the share is not cut — there is nothing to cut
// it against — so its column is simply taller than the rest.
func (c *CSS) multiBalance(b *multiBatch) int {
	total := 0
	for i := range b.items {
		total += b.items[i].mt + b.items[i].h + b.items[i].mb
	}
	if b.count < 1 || total == 0 {
		return 0
	}
	return max((total+b.count-1)/b.count, 1)
}

// solveColumns reads the count and the width apart. An explicit column-count
// says how many; a column-width says how many of them fit, and the smaller of
// the two is what the container makes. With neither there is one column and
// the block flows as it would have without any of this. The width is measured
// the way every other length in the sheet is — through the template's length,
// so it answers to the window's unit like the gutter beside it.
func (b *multiBatch) solveColumns(st css.Style) {
	n := 0
	if st.ColumnCount > 0 {
		n = st.ColumnCount
	}
	if st.Has("column-width") && gridLengthDefinite(st.ColumnWidth) {
		w := b.css.style.length(st, st.ColumnWidth, b.contentW)
		if w > 0 {
			if fit := (b.contentW + b.gap) / (w + b.gap); n == 0 || fit < n {
				n = fit
			}
		}
	}
	if n < 1 {
		n = 1
	}
	b.count = n
	if n == 1 {
		b.colW = b.contentW
		return
	}
	b.colW = max((b.contentW-b.gap*(n-1))/n, 0)
}

// flow cuts the items into the columns and remembers how many draw passes the
// pieces need. Everything it writes is derived from the box the container was
// given and the measured items, so the second call a nested container makes —
// at the box its parent really gave it — answers for that box.
func (b *multiBatch) flow() {
	b.solveColumns(b.st)
	b.colH = 0
	for i := range b.items {
		it := &b.items[i]
		it.w, it.h = b.css.multiItemBox(it.st, it.nw, it.nh, b.colW)
		it.pieces = it.pieces[:0]
	}

	// Only a container whose height was decided has one to cut a box against;
	// the rest balance their content over the columns they were given.
	b.cuts = b.definite && b.contentH > 0
	if b.cuts {
		b.room = b.css.multiColumnRoom(b, b.st)
	} else {
		b.room = b.css.multiBalance(b)
	}

	col, y, prevMB := 0, 0, 0
	for i := range b.items {
		it := &b.items[i]
		for {
			gap := max(prevMB, it.mt)
			free := max(b.room-y-gap, 0)
			last := col >= b.count-1
			switch {
			case it.h <= free:
				it.pieces = append(it.pieces, multiPiece{col: col, y: y + gap, h: it.h})
				y += gap + it.h
				prevMB = it.mb
			case !b.cuts && !last && y > 0:
				// A balancing container starts the next column and puts the
				// whole box in it: the share is a guide, not a wall. A column
				// with nothing in it yet takes the box even so, because a
				// break before the first box would only move the tall column
				// to the other side.
				col++
				y, prevMB = 0, 0
				continue
			case b.cuts && !last && !it.keep:
				// A container with a height cuts the box at the edge of the
				// column and continues it in the next one.
				if len(it.pieces) >= multiPieces {
					it.pieces = append(it.pieces, multiPiece{col: col, y: y + gap, h: it.h})
					y += gap + it.h
					prevMB = it.mb
					break
				}
				it.pieces = append(it.pieces, multiPiece{col: col, y: y + gap, h: free})
				it.h -= free
				col++
				y, prevMB = 0, 0
				continue
			case b.cuts && !last:
				// A box that asked not to be cut starts the next column whole
				// instead, which is what break-inside: avoid buys it.
				col++
				y, prevMB = 0, 0
				continue
			case b.cuts:
				// The last column: the box is cut at the edge and the rest of
				// it hangs below the columns, which is where content that
				// does not fit the columns a container asked for goes.
				if free > 0 {
					it.pieces = append(it.pieces, multiPiece{col: col, y: y + gap, h: free})
					y += free
				}
				it.pieces = append(it.pieces, multiPiece{col: col, y: y, h: it.h - free})
				y += it.h - free
				prevMB = it.mb
			default:
				// A balancing container whose column is out of room, or empty,
				// or the last one: the box goes in whole and that column is the
				// taller one.
				it.pieces = append(it.pieces, multiPiece{col: col, y: y + gap, h: it.h})
				y += gap + it.h
				prevMB = it.mb
			}
			// The container is as tall as its tallest column, which is not
			// always the one the flow finished in.
			b.colH = max(b.colH, y+prevMB)
			break
		}
	}

	// A container that declared no height of its own is as tall as the columns
	// came out, so a parent that asks for the box it wants has one to place.
	if !b.cuts {
		b.h = max(b.colH+b.top+b.bottom, 0)
		b.contentH = b.h - b.top - b.bottom
	}

	// One pass per fragment: every widget is drawn once for each piece it owns.
	passes := 1
	for i := range b.items {
		passes = max(passes, len(b.items[i].pieces))
	}
	b.passes = passes
}
