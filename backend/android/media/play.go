//go:build android

package media

import "time"

// Play starts, or carries on after a pause.
func (p *Player) Play() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.playing {
		return
	}
	p.playing = true
	p.started = time.Now()
}
