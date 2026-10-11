package canvas

// PresentChanges commits the frame: it measures the rows that were written
// against the baseline, copies the pixels that really moved into it, and
// answers the rectangle they cover — the only part of the frame the display
// has to be sent. A frame that changed nothing answers an empty area and
// writes nothing.
//
// Measuring here rather than at each write is what makes the answer exact.
// Drawing is overdraw: a card paints a shadow, then a background over it, then
// a border over that, and a frame later the same three land on the same
// pixels. Every one of those writes looked like a change when it happened and
// none of them was one by the end. The canvas ends up holding what was drawn
// either way; what the screen needs is only known once the drawing is over.
func (cv *Canvas) PresentChanges() Area {
	if cv.base == nil {
		return Area{}
	}
	area := Area{}
	for _, y := range cv.rows {
		s := &cv.spans[y]
		left, right, moved := changedSpan(
			cv.Pixels[y*cv.Stride+s.left:y*cv.Stride+s.right+1],
			cv.base[y*cv.Width+s.left:y*cv.Width+s.right+1])
		if moved {
			first, last := s.left+left, s.left+right
			copy(cv.base[y*cv.Width+first:y*cv.Width+last+1],
				cv.Pixels[y*cv.Stride+first:y*cv.Stride+last+1])
			area = growArea(area, first, last, y)
		}
		*s = rowSpan{left: -1}
	}
	cv.rows = cv.rows[:0]
	cv.change = area
	return area
}
