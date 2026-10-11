package template

import "github.com/gabrielluizsf/antui/template/css"

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
