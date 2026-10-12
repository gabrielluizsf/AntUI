//go:build android

package media

// Playing reports whether the clock is running.
func (p *Player) Playing() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.playing
}
