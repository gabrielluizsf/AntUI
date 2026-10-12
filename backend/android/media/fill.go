//go:build android

package media

// fill hands the device the next block, decoding more when the buffer runs
// low. It runs on the audio device's own goroutine.
func (p *Player) fill(block []float32) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i := range block {
		block[i] = 0
	}
	if p.closed || !p.playing {
		return
	}
	// Top the buffer up. A decode gives a few thousand samples at a time,
	// so a handful of goes covers any block a device will ask for.
	for p.buffered() < len(block) && !p.audio.Ended() {
		if !p.decodeSound() {
			break
		}
	}
	n := min(p.buffered(), len(block))
	for i := range n {
		block[i] = p.pcm[(p.pcmHead+i)%len(p.pcm)]
	}
	p.pcmHead = (p.pcmHead + n) % len(p.pcm)
	p.played += int64(n)
	if p.audio.Ended() && p.buffered() == 0 {
		p.ended = true
	}
}
