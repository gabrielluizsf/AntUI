package audio

// Fill is what [antui.OpenAudio] calls: it overwrites the given interleaved
// stereo block with whatever is sounding. It is safe to call from wherever
// the device calls it.
func (p *Player) Fill(block []float32) {
	if p == nil {
		for i := range block {
			block[i] = 0
		}
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	for i := 0; i+1 < len(block); i += 2 {
		var l, r float32
		for s := range p.slots {
			v := p.read(&p.slots[s])
			if v == 0 {
				continue
			}
			// A slight balance by pan; sounds default to the middle.
			l += v * (1.0 - p.slots[s].pan) * 0.5
			r += v * (1.0 + p.slots[s].pan) * 0.5
		}
		block[i], block[i+1] = clamp(l), clamp(r)
	}
}
