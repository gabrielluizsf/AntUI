//go:build android

package media

// Ended reports whether the file has played out.
func (p *Player) Ended() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.ended
}
