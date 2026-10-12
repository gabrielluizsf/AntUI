//go:build android

package media

import (
	"errors"
	"time"

	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

// Seek jumps to a position. Both tracks jump to the nearest key frame at or
// before it, which for the picture may be as much as a second earlier —
// there is no other kind of seek in a compressed stream.
func (p *Player) Seek(to time.Duration) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return errors.New("media: the player is closed")
	}
	for _, d := range []*ndk.Decoder{p.video, p.audio} {
		if d == nil {
			continue
		}
		if err := d.Seek(to); err != nil {
			return err
		}
	}
	p.base = to
	p.started = time.Now()
	p.played = 0
	p.pcmHead, p.pcmTail = 0, 0
	p.haveNext = false
	p.curPts, p.nextPts = to, to
	p.ended = false
	return nil
}
