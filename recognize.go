package antui

// recognize turns this frame's touches into gestures. It runs from Begin,
// after the backend has pushed everything, so it sees a complete frame.
//
// The time is passed in rather than read here, because half of what a
// gesture is is a duration — a tap is a press that was short and a long
// press is one that was not — and a recogniser that reads the clock itself
// can only be tested by waiting.
func (win *Window) recognize(now float64) {
	g := &win.gesture
	g.out = Gestures{}

	win.recognizeOne(now)
	win.recognizeTwo()
}
