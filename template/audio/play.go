package audio

// play puts a sound on a free voice, or on the voice closest to finishing.
// pitch speeds the sound up or down so one sample can say many things.
func (p *Player) play(s *Sound, pitch float64, pan float32) {
	if p == nil || s == nil || len(s.sample) == 0 {
		return
	}
	best := -1
	bestLeft := 0.0
	for i := range p.slots {
		sl := &p.slots[i]
		if len(sl.sample) == 0 || sl.pos >= float64(len(sl.sample))-1 {
			best = i
			break
		}
		left := float64(len(sl.sample)) - sl.pos
		if best < 0 || left < bestLeft {
			best, bestLeft = i, left
		}
	}
	sl := &p.slots[best]
	sl.sample = s.sample
	sl.rate = s.rate
	sl.pos = 0
	sl.pitch = pitch
	sl.gain = s.gain
	sl.pan = pan
}
