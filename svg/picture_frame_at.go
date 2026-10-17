package svg

import (
	"math"

	"github.com/gabrielluizsf/antui/canvas"
)

// indexAt is the frame showing at a moment on the clock: the frames laid end to
// end into one loop, the moment wrapped into it, and the frame whose stretch of
// it holds the moment. A loop of no time at all keeps the first frame — there
// is nowhere for the moment to advance to.
func (f *pictureFrames) indexAt(now float64) int {
	if f == nil || len(f.frames) == 0 {
		return 0
	}
	var total float64
	for _, s := range f.show {
		total += s
	}
	if total <= 0 {
		return 0
	}
	t := math.Mod(now, total)
	if t < 0 {
		t += total
	}
	for i, s := range f.show {
		if t < s {
			return i
		}
		t -= s
	}
	return len(f.frames) - 1
}

// at is the frame showing at that moment, or nothing when there is no picture
// to show — an element whose bytes have yet to arrive draws its hole.
func (f *pictureFrames) at(now float64) *canvas.Canvas {
	if f == nil || len(f.frames) == 0 {
		return nil
	}
	return f.frames[f.indexAt(now)]
}
