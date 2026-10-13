//go:build android

package media

import "time"

// Position is where the file is up to.
func (p *Player) Position() time.Duration {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.position()
}
