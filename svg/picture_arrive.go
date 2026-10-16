package svg

import "github.com/gabrielluizsf/antui/canvas"

// fitted is the frame the picture is showing at a moment on the clock, with
// where it goes — the box read in the same breath as the frame, because a side
// of the box the drawing left open is said by the picture itself when it
// arrives, and a paint walking the drawing while a fetch off the frame is
// landing has to see the two together. Nothing is there until the bytes an
// `<image>` named have arrived.
func (p *Picture) fitted(now float64) (*canvas.Canvas, pictureFit) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.frames.at(now), pictureFit{box: p.Box, fit: p.Fit, align: p.Align}
}

// fittedBox is where the picture goes, read through the lock for the reason
// [Picture.fitted] reads it with the frame: a measure of what a drawing paints
// over may run while a picture is landing, and a box of nothing is the honest
// answer for a picture that is not there yet.
func (p *Picture) fittedBox() [4]float64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.Box
}

// frameIndex is which frame of the loop that moment holds — the number the
// clock picks by — read for the drawing's animation tick, which says whether a
// canvas painted a moment ago still shows what is on the clock.
func (p *Picture) frameIndex(now float64) int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.frames.indexAt(now)
}

// arrive puts the fetched picture into the slot the node was built with, so
// the next paint of the drawing shows it without the drawing being read again
// — and finishes the box while it is at it, taking the sides the drawing left
// open from the size the picture answered with.
func (p *Picture) arrive(pic *pictureFrames) {
	ownW, ownH := pic.ownSize()
	p.mu.Lock()
	p.Box = autoBox(p.declared, p.autoW, p.autoH, ownW, ownH)
	p.frames = pic
	p.mu.Unlock()
}
