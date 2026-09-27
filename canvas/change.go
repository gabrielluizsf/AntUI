package canvas

import "slices"

// A canvas can be given a baseline: the pixels of the frame that was last
// presented to the screen. Every write from then on is remembered, and at the
// end of the frame the remembered rows are measured against the baseline, so
// a program that redraws the same picture every frame asks the display for
// nothing at all, and one that moves a cursor across it asks for six pixels
// rather than the whole window.
//
// The measuring is where the work is bounded: only rows that were written are
// looked at, and only across the span that was written. Nothing ever walks the
// whole canvas, and a frame that wrote the whole canvas pays one pass of
// memcmp and sends nothing.
//
//	win.Canvas().Compare(previousFrame)   // at the start of the frame
//	...
//	dirty := win.Canvas().PresentChanges() // at the end of it
type rowSpan struct {
	left  int // the first pixel written on this row; -1 when none was
	right int // and the last
}

// Compare makes this canvas measure its writes against base, the pixels of the
// frame that was last presented, and forgets whatever it had recorded. A
// baseline that is not this canvas's own size turns the measuring off: a
// canvas that was resized has nothing to compare against, and everything it
// draws is new.
func (cv *Canvas) Compare(base []Color) {
	if cv.Pixels == nil || len(base) != cv.Width*cv.Height {
		cv.Forget()
		return
	}
	cv.base = base
	if cap(cv.spans) < cv.Height {
		cv.spans = make([]rowSpan, cv.Height)
	}
	cv.spans = cv.spans[:cv.Height]
	for i := range cv.spans {
		cv.spans[i] = rowSpan{left: -1}
	}
	cv.rows = cv.rows[:0]
	cv.change = Area{}
}

// Forget drops the baseline, so writes go straight to the pixels again with
// nothing recorded. The Dirty bounds still stand: they say what was written,
// which is a different question from what changed.
func (cv *Canvas) Forget() {
	cv.base = nil
	cv.rows = cv.rows[:0]
	cv.change = Area{}
}

// Comparing reports whether this canvas is measuring its writes against a
// baseline.
func (cv *Canvas) Comparing() bool { return cv.base != nil }

// Changed is the rectangle that covered every pixel that really moved when the
// last frame was committed by [Canvas.PresentChanges]: the part of the screen
// that had to be sent for it. It is empty before the first commit, and empty
// again after one that changed nothing.
func (cv *Canvas) Changed() (Area, bool) {
	return cv.change, cv.change.Width > 0
}

// Changes walks the rectangles that really moved in the last committed frame.
// The window uploads their bounding box; a backend that damages per box —
// Wayland, Android — walks this instead and sends only the pieces that moved.
func (cv *Canvas) Changes(yield func(area Area) bool) {
	if cv.change.Width <= 0 {
		return
	}
	x, y := cv.change.X, cv.change.Y
	w, h := cv.change.Width, cv.change.Height
	for dy := range h {
		if !yield(Area{X: x, Y: y + dy, Width: w, Height: 1}) {
			return
		}
	}
}

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

// growArea is the bounding box of a rectangle and a run of pixels on a row.
func growArea(a Area, x0, x1, y int) Area {
	if a.Width <= 0 {
		return Area{X: x0, Y: y, Width: x1 - x0 + 1, Height: 1}
	}
	if x0 < a.X {
		a.Width += a.X - x0
		a.X = x0
	}
	if x1 >= a.X+a.Width {
		a.Width = x1 - a.X + 1
	}
	if y < a.Y {
		a.Height += a.Y - y
		a.Y = y
	}
	if y >= a.Y+a.Height {
		a.Height = y - a.Y + 1
	}
	return a
}

// markSpan remembers that a row was written between x0 and x1. The span is
// the widest the row was written, which is all the measuring needs: the run is
// narrowed to the pixels that really moved when the frame is committed.
func (cv *Canvas) markSpan(y, x0, x1 int) {
	if cv.base == nil || x1 < x0 || y < 0 || y >= len(cv.spans) {
		return
	}
	s := &cv.spans[y]
	if s.left < 0 {
		s.left, s.right = x0, x1
		cv.rows = append(cv.rows, y)
		return
	}
	s.left, s.right = min(s.left, x0), max(s.right, x1)
}

// markRun is markSpan for a path that wrote a run of pixels between x0 and
// x1 without knowing what each of them became — a blend, a copy, a filter.
// The run is half-open, as the write loops that call it are.
func (cv *Canvas) markRun(y, x0, x1 int) {
	if x1 > x0 {
		cv.markSpan(y, x0, x1-1)
	}
}

// markRegion is markRun for a whole rectangle, which is what an in-place pass
// — a blur, a filter, a fade — hands over: it rewrote every pixel of it.
func (cv *Canvas) markRegion(r Area) {
	if cv.base == nil {
		return
	}
	x0, y0 := max(r.X, 0), max(r.Y, 0)
	x1, y1 := min(r.X+r.Width, cv.Width), min(r.Y+r.Height, cv.Height)
	for y := y0; y < y1; y++ {
		cv.markRun(y, x0, x1)
	}
}

// changedSpan is the first and the last pixel of two equal-length runs that
// differ, and whether they differ at all. The whole-run comparison is one
// memcmp; the narrowing walk only runs on a run that moved, and then usually
// stops at the first pixel.
func changedSpan(a, b []Color) (int, int, bool) {
	n := min(len(a), len(b))
	if slices.Equal(a[:n], b[:n]) {
		return 0, 0, false
	}
	left := 0
	for a[left] == b[left] {
		left++
	}
	right := n - 1
	for a[right] == b[right] {
		right--
	}
	return left, right, true
}
