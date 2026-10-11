package canvas

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
