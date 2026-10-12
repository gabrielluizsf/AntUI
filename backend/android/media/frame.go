//go:build android

package media

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// Frame is the picture that belongs on screen now, and whether there is one.
//
// The canvas belongs to the player and is handed out again and again; a
// frame that has to outlive the next call has to be copied. It reports false
// before the first frame has been decoded and after the file has ended.
func (p *Player) Frame() (*canvas.Canvas, bool, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed || p.video == nil {
		return nil, false, nil
	}
	now := p.position()

	// Decode ahead until the frame in hand is the right one for now. More
	// than a few in one go means the decoder is behind, and running the loop
	// to the end would make it later still — so it catches up a little each
	// frame instead of stalling one.
	for range 4 {
		if !p.haveNext {
			if err := p.decodeNext(); err != nil {
				return nil, false, err
			}
		}
		if !p.haveNext || p.nextPts > now {
			break
		}
		p.cur, p.next = p.next, p.cur
		p.curPts = p.nextPts
		p.haveNext = false
	}
	if p.curPts == 0 && !p.started.IsZero() && p.cur == nil {
		return nil, false, nil
	}
	return p.cur, true, nil
}
