package audio

// read advances a voice and returns its next sample, applying a short fade so
// a sound never starts or ends on a hard edge.
func (p *Player) read(s *slot) float32 {
	if s == nil || len(s.sample) == 0 {
		return 0
	}
	if s.pos >= float64(len(s.sample))-1 {
		s.sample = nil
		return 0
	}
	i := int(s.pos)
	frac := float32(s.pos - float64(i))
	v := s.sample[i] + (s.sample[i+1]-s.sample[i])*frac

	// Attack over three milliseconds, release over four, counted in the
	// source's time so a slowed-down sound does not click louder.
	srcPerOut := float64(s.rate) * s.pitch / float64(p.Rate)
	attack := float64(s.rate) * 0.003
	release := float64(s.rate) * 0.004
	env := float32(minf(float64(s.pos)/attack, 1))
	remaining := float64(len(s.sample)-1) - s.pos
	if rem := remaining / release; rem < 1 {
		env *= float32(maxf(rem, 0))
	}

	s.pos += srcPerOut
	return v * env * s.gain
}
