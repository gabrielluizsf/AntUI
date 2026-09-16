package audio

import "sync"

// A Player mixes the short sounds into the stereo stream [antui.OpenAudio]
// asks for. Sounds are triggered from the game's own goroutine and read from
// the device's goroutine, so everything the two share goes behind the lock:
// the waits are never long enough to be heard.
//
// A handful of voices — eight, enough for even rapid typing on top of clicks —
// play at once. When every voice is busy the one nearest its end is taken
// over, so a fresh sound is never dropped for an old one.
type Player struct {
	mu    sync.Mutex
	Rate  int // the rate Fill is called at, in frames a second
	slots [8]slot
}

// A slot is one sound, in the middle of being played.
type slot struct {
	sample []float32
	rate   int     // the sound's own sample rate
	pos    float64 // where in the sample, in source samples
	pitch  float64 // how fast to go through it (1 = as recorded)
	gain   float32
	pan    float32 // -1 to 1, the channel balance
}

// NewPlayer makes a mixer that writes to a stream running at the given rate.
// Hand its Fill to [antui.OpenAudio]; the rate must be the one the stream
// really runs at, which is what [antui.Audio.Rate] reports.
//
//	player := audio.NewPlayer(stream.Rate())
//	stream, err := antui.OpenAudio("game", 44100, player.Fill)
func NewPlayer(rate int) *Player {
	if rate <= 0 {
		rate = 44100
	}
	return &Player{Rate: rate}
}

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

func clamp(v float32) float32 {
	if v > 1 {
		return 1
	}
	if v < -1 {
		return -1
	}
	return v
}

func minf(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func maxf(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
