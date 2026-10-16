package svg

import "github.com/gabrielluizsf/antui/canvas"

// pixels is the picture as it is right now, which is nothing until the bytes an
// `<image>` named have arrived. Reading it through the lock is what lets a
// fetch running off the frame put the picture in while a paint is walking the
// drawing: the field is written once, and every read after that sees it.
func (p *Picture) pixels() *canvas.Canvas {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.Pixels
}

// arrive puts the fetched picture into the slot the node was built with, so the
// next paint of the drawing shows it without the drawing being read again.
func (p *Picture) arrive(pic *canvas.Canvas) {
	p.mu.Lock()
	p.Pixels = pic
	p.mu.Unlock()
}
