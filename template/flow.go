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
		col, y, prevMB = b.placeItem(&b.items[i], col, y, prevMB)
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
