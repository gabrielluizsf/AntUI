//go:build android

package media

// Pause stops the clock. Nothing is thrown away.
func (p *Player) Pause() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.playing {
		return
	}
	p.base = p.position()
	p.playing = false
}
