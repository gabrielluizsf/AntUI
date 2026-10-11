package template

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
