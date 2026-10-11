package template

import "github.com/gabrielluizsf/antui/template/css"

// alignCross plants every item across the cross axis of its line, offset by
// where the line itself starts: stretch grows an auto-sized item to fill, and
// the other keywords park it against a line edge or in the middle.
func (b *flexBatch) alignCross(line []int, cross, off int) {
	for _, idx := range line {
		it := &b.items[idx]
		a := b.alignItems()
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
