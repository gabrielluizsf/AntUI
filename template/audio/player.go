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
