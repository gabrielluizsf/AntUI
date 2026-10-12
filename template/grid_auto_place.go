package template

import "github.com/gabrielluizsf/antui/template/css"

// gridAutoPlace runs the automatic placement passes: the items locked to one
// line come first, in the order they were added, and everything else follows
// in that same order, so an item with one definite axis takes the cursor where
// the items before it left it.
func (b *gridBatch) gridAutoPlace(columns, rows *[]css.GridTrack, occupied map[[2]int]bool) {
	order := b.gridOrder()
	rowCursor, colCursor := 0, 0
	place := func(index int) {
		it := &b.items[index]
		if !it.out && !(it.colDefinite && it.rowDefinite) {
			b.placeAutomatic(it, columns, rows, occupied, &rowCursor, &colCursor)
		}
	}
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
		return
	}
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
