package canvas

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
