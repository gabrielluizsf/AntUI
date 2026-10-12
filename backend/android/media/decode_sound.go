//go:build android

package media

// decodeSound puts one decoded block into the ring, and reports whether it
// got anything.
func (p *Player) decodeSound() bool {
	room := len(p.pcm) - p.buffered() - 1
	if room <= 0 {
		return false
	}
	scratch := make([]float32, min(room, 8192))
	n, _, err := p.audio.Samples(scratch)
	if err != nil || n == 0 {
		return false
	}
	for i := range n {
		p.pcm[p.pcmTail] = scratch[i]
		p.pcmTail = (p.pcmTail + 1) % len(p.pcm)
	}
	return true
}
