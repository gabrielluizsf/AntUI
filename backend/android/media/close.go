//go:build android

package media

// Close stops everything and gives it back.
func (p *Player) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return nil
	}
	p.closed = true
	p.playing = false
	if p.video != nil {
		p.video.Close()
	}
	if p.audio != nil {
		p.audio.Close()
	}
	p.cur, p.next = nil, nil
	return nil
}
