package audio

import "math"

// A Sound is one sample a template can play: a short recording, decoded once
// at load time and replayed as often as needed. Playing needs a Player that
// has been [Use]d, or nothing at all happens — a game on a silent machine
// keeps running.
type Sound struct {
	sample []float32
	rate   int
	gain   float32
}

// Play sounds the sample as recorded.
func (s *Sound) Play() { s.play(1, 0) }

// PlayStep sounds the sample at a pitch stepping up a gentle scale with the
// given step, which is how a field of text can be typed at without every key
// sounding identical. Free of charge for any sound, so sliders may step up
// and down with it too.
func (s *Sound) PlayStep(step int) {
	s.play(pitchFor(step), 0)
}

// PlayTuned sounds the sample at a given playback speed. 1 is the recorded
// pitch, 2 an octave up, 0.5 an octave down.
func (s *Sound) PlayTuned(pitch float64) { s.play(pitch, 0) }

// play is the one door everything goes through: to the default device, if
// there is one.
func (s *Sound) play(pitch float64, pan float32) {
	if s == nil {
		return
	}
	if dev := defaultPlayer.Load(); dev != nil {
		dev.play(s, pitch, pan)
	}
}

// pitchFor maps a keystroke's index onto a short rising scale, so tying a
// pitch to where the cursor sits makes typing climb rather than repeat. A
// whole octave of semitones keeps it recognisable whatever sound is used.
func pitchFor(step int) float64 {
	// The indices are of semitones in a major scale, two octaves of it.
	scale := [16]int{0, 2, 4, 5, 7, 9, 11, 12, 14, 16, 17, 19, 21, 23, 24, 26}
	if step < 0 {
		step = -step
	}
	semi := scale[step%len(scale)]
	return math.Pow(2, float64(semi)/12)
}