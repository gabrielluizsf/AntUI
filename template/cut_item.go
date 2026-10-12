package template

// cutItem places a box that overflows a column of a container with a decided
// height. It returns the new cursor and whether the flow must retry the item
// in the column it moved to.
func (b *multiBatch) cutItem(it *multiItem, col, y, gap, free int, last bool) (int, int, int, bool) {
	if !last && !it.keep {
		// A container with a height cuts the box at the edge of the column
		// and continues it in the next one.
		if len(it.pieces) >= multiPieces {
			it.pieces = append(it.pieces, multiPiece{col: col, y: y + gap, h: it.h})
			return col, y + gap + it.h, it.mb, false
		}
		it.pieces = append(it.pieces, multiPiece{col: col, y: y + gap, h: free})
		it.h -= free
		return col + 1, 0, 0, true
	}
	if !last {
		// A box that asked not to be cut starts the next column whole
		// instead, which is what break-inside: avoid buys it.
		return col + 1, 0, 0, true
	}
	// The last column: the box is cut at the edge and the rest of it hangs
	// below the columns, which is where content that does not fit the columns
	// a container asked for goes.
	if free > 0 {
		it.pieces = append(it.pieces, multiPiece{col: col, y: y + gap, h: free})
		y += free
	}
	it.pieces = append(it.pieces, multiPiece{col: col, y: y, h: it.h - free})
	return col, y + it.h - free, it.mb, false
}
