package template

// placeItem puts one measured item into the columns, starting at col/y with
// prevMB the margin of what came before, and returns the flow cursor for the
// next item. A box too tall for the column is cut into pieces, or moved whole
// to the next one, per the container's rules.
func (b *multiBatch) placeItem(it *multiItem, col, y, prevMB int) (int, int, int) {
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
			// A balancing container starts the next column and puts the whole
			// box in it: the share is a guide, not a wall.
			col++
			y, prevMB = 0, 0
			continue
		case b.cuts:
			var retry bool
			col, y, prevMB, retry = b.cutItem(it, col, y, gap, free, last)
			if retry {
				continue
			}
		default:
			// The box goes in whole and this column becomes the taller one.
			it.pieces = append(it.pieces, multiPiece{col: col, y: y + gap, h: it.h})
			y += gap + it.h
			prevMB = it.mb
		}
		// The container is as tall as its tallest column, which is not
		// always the one the flow finished in.
		b.colH = max(b.colH, y+prevMB)
		return col, y, prevMB
	}
}
