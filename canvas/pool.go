package canvas

import "sync"

// layerPool hands out scratch canvases: the copies a shadow, a blur or a
// backdrop is made in before it is composited. A shadow is the same blurred
// box drawn again on the next frame, and every one of them used to make a
// canvas of its own — a box's worth of pixels, and inside it, a box's worth of
// blur scratch, thrown away the moment it was composited. The pool keeps one
// of each size around instead, which is the whole difference between a frame
// that allocates nothing and a frame that hands the collector a few megabytes
// every time a card on the screen has a shadow under it.
var layerPool sync.Pool

// acquireLayer returns a cleared scratch canvas exactly w by h. Its pixels
// are all transparent, its clip is the whole of it, and anything bigger it was
// last used as is trimmed away, so a caller can treat it as if it were new.
func acquireLayer(w, h int) (*Canvas, bool) {
	if w <= 0 || h <= 0 {
		return nil, false
	}
	cv, _ := layerPool.Get().(*Canvas)
	if cv == nil {
		cv = &Canvas{format: ARGB32}
	}
	n := w * h
	if cap(cv.Pixels) < n {
		cv.Pixels = make([]Color, n)
	}
	cv.Pixels = cv.Pixels[:n]
	clear(cv.Pixels)
	cv.Width, cv.Height, cv.Stride = w, h, w
	cv.Clip = Area{0, 0, w, h}
	cv.Dirty = Area{}
	cv.base, cv.spans, cv.rows, cv.change = nil, nil, nil, Area{}
	return cv, true
}

// releaseLayer hands a scratch canvas back. Anything holding on to it past
// this point is drawing into a canvas someone else is about to use.
func releaseLayer(cv *Canvas) {
	if cv == nil {
		return
	}
	// A scratch canvas can carry a blur's float buffers, and those are worth
	// more than the pixels they sit on; a shadow the size of a card is a
	// blur the size of a card, and letting go of it means paying for the next
	// one in full.
	layerPool.Put(cv)
}
