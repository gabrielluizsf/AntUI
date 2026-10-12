//go:build android

package media

import (
	"unsafe"

	"github.com/gabrielluizsf/antui/canvas"
)

// decodeNext pulls one frame out of the decoder into the spare canvas.
func (p *Player) decodeNext() error {
	if p.next == nil {
		return nil
	}
	words := unsafe.Slice((*uint32)(unsafe.Pointer(&p.next.Pixels[0])), len(p.next.Pixels))
	w, h, pts, ok, err := p.video.Frame(words, p.next.Stride)
	if err != nil {
		return err
	}
	if p.video.Ended() {
		p.ended = true
	}
	if !ok {
		return nil
	}
	if w != p.next.Width || h != p.next.Height {
		// The decoder settled on a different size than the container
		// claimed. Make canvases that fit and lose this one frame.
		cv, err := canvas.NewCanvas(w, h)
		if err != nil {
			return err
		}
		p.next = cv
		if spare, err := canvas.NewCanvas(w, h); err == nil {
			p.cur = spare
		}
		p.info.Width, p.info.Height = w, h
		return nil
	}
	p.nextPts = pts
	p.haveNext = true
	return nil
}
