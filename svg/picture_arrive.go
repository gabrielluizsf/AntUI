package svg

import "github.com/gabrielluizsf/antui/canvas"

// at is the frame the picture is showing at a moment on the clock, which is
// nothing until the bytes an `<image>` named have arrived. Reading it through
// the lock is what lets a fetch running off the frame put the picture in while
// a paint is walking the drawing: the slot is written once, and every read
// after that sees it.
func (p *Picture) at(now float64) *canvas.Canvas {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.frames.at(now)
}

// frameIndex is which frame of the loop that moment holds — the number
// [Picture.at] picks by — read for the drawing's animation tick, which says
// whether a canvas painted a moment ago still shows what is on the clock.
func (p *Picture) frameIndex(now float64) int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.frames.indexAt(now)
}

// arrive puts the fetched picture into the slot the node was built with, so the
// next paint of the drawing shows it without the drawing being read again.
func (p *Picture) arrive(pic *pictureFrames) {
	p.mu.Lock()
	p.frames = pic
	p.mu.Unlock()
}
