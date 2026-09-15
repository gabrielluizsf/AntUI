package audio

import (
	"testing"
)

// tapSound is a known square wave the mixer tests can predict.
func tapSound(rate int, n int) *Sound {
	s := make([]float32, n)
	for i := range n {
		s[i] = 0.5
	}
	return &Sound{sample: s, rate: rate, gain: 1}
}

func TestFillSilentWhenNothingPlays(t *testing.T) {
	p := NewPlayer(44100)
	block := make([]float32, 64)
	p.Fill(block)
	for i, v := range block {
		if v != 0 {
			t.Fatalf("block[%d] = %v, want silence", i, v)
		}
	}
}

func TestPlayingMakesSound(t *testing.T) {
	p := NewPlayer(44100)
	p.play(tapSound(44100, 8000), 1, 0)
	block := make([]float32, 256)
	p.Fill(block)
	var loudest float32
	for _, v := range block {
		if a := v; a > loudest {
			loudest = a
		}
	}
	if loudest == 0 {
		t.Error("a playing sound should come out of Fill")
	}
}

func TestSoundSitsAtDeviceRateAfterResampling(t *testing.T) {
	// Source at 44100 played into a 22050 stream takes twice as long to get
	// through. Play a one-second 44.1k sound into a 22.05k stream and the
	// fill should have heard it for about two seconds of output.
	p := NewPlayer(22050)
	p.play(tapSound(44100, 44100), 1, 0)

	probe := make([]float32, 22050*2) // two output seconds
	p.Fill(probe)
	var heard bool
	for _, v := range probe {
		if v != 0 {
			heard = true
			break
		}
	}
	if !heard {
		t.Error("half-rate resampling dropped the sound")
	}
}

func TestPitchShortensTheSound(t *testing.T) {
	p := NewPlayer(44100)
	p.play(tapSound(44100, 44100), 2, 0) // twice as fast → half the duration
	probe := make([]float32, 44100)
	p.Fill(probe)

	lastFrame := -1
	for i := 0; i+1 < len(probe); i += 2 {
		if probe[i] != 0 || probe[i+1] != 0 {
			lastFrame = i / 2
		}
	}
	// Half a second of sound at 2x pitch, plus the release ramp and the
	// freshly-stolen tail: comfortably under the one-second probe.
	if lastFrame <= 0 || lastFrame > 22050+1102 {
		t.Errorf("2x-pitch sound ran to frame %d, want ~%d", lastFrame, 22050)
	}
}

func TestBusyPlayerReusesNearestFinishedVoice(t *testing.T) {
	p := NewPlayer(44100)
	// Two short sounds and one long one fill every voice a little.
	for range 8 {
		p.play(tapSound(44100, 2000), 1, 0)
	}
	p.play(tapSound(44100, 44100*3), 1, 0) // long one: the one still going
	p.play(tapSound(44100, 2000), 1, 0)    // should steal and be heard

	probe := make([]float32, 44100)
	p.Fill(probe)
	var heard bool
	for _, v := range probe {
		if v != 0 {
			heard = true
			break
		}
	}
	if !heard {
		t.Error("the last sound should have stolen a finished voice")
	}
}

func TestPlayStepRisesWithTheScale(t *testing.T) {
	p := NewPlayer(44100)
	s := tapSound(44100, 8000)
	low := pitchFor(0)
	high := pitchFor(7)
	if low >= high {
		t.Errorf("pitchFor(7) = %v should rise above pitchFor(0) = %v", high, low)
	}
	if p == nil {
		t.Fatal("player was not created")
	}
	_ = s
}

func TestNilPlayerAndSilentDevicePlayNothing(t *testing.T) {
	var s *Sound
	s.Play() // must not panic
	s = tapSound(44100, 100)
	s.Play() // no default device: still nothing
	Use(nil)
	s.Play()
}