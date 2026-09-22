//go:build android

package media

import (
	"errors"
	"sync"
	"time"
	"unsafe"

	"github.com/gabrielluizsf/antui/backend/android/ndk"
	"github.com/gabrielluizsf/antui/canvas"
)

// Info is what a file holds.
type Info struct {
	Duration time.Duration
	// Width and Height are the picture's, after the rotation the file
	// carries has been applied — which is the size it should be drawn at.
	Width, Height int
	// Rotation is how far the picture is stored turned. A video shot on a
	// phone held sideways is stored as if it were not, with a 90 here.
	Rotation int
	// SampleRate and Channels are the sound's, and are zero when the file
	// has none.
	SampleRate, Channels int
	// HasVideo and HasAudio say which tracks are there.
	HasVideo, HasAudio bool
}

// Player plays a file: the picture into a [canvas.Canvas], the sound through
// whatever the caller hands the samples to.
//
// It does not draw and it does not play. It decodes, keeps time, and hands
// over the frame whose moment has come — everything else is the app's, which
// is what lets a video be a texture on something, or half the screen, or
// paused while the rest of the game carries on.
type Player struct {
	info Info

	video *ndk.Decoder
	audio *ndk.Decoder

	mu sync.Mutex
	// cur is the frame being shown and next the one after it, decoded ahead
	// so that the moment it is due there is no wait. Two is enough: the
	// point is to be one ahead, not to buffer.
	cur, next       *canvas.Canvas
	curPts, nextPts time.Duration
	haveNext        bool

	// The clock. When there is sound it is the sound that keeps time —
	// counting frames actually handed to the device — because an ear notices
	// a hundredth of a second and an eye does not. With no sound it is the
	// wall clock.
	played  int64
	started time.Time
	base    time.Duration
	playing bool

	pcm     []float32
	pcmHead int
	pcmTail int

	ended  bool
	closed bool
}

// Open reads a file and gets both its tracks ready.
//
// The path is a file on this device. An address from
// [antui/backend/android/picker] is not one, and opening it needs the content
// resolver
func Open(path string) (*Player, error) {
	m, err := ndk.OpenMedia(path)
	if err != nil {
		return nil, err
	}
	p := &Player{}

	if track, ok := m.Track(true); ok {
		d, err := m.DecodeVideo(track, track.Rotation)
		if err != nil {
			return nil, err
		}
		p.video = d
		p.info.HasVideo = true
		p.info.Rotation = track.Rotation
		p.info.Width, p.info.Height = track.Width, track.Height
		if track.Rotation == 90 || track.Rotation == 270 {
			p.info.Width, p.info.Height = track.Height, track.Width
		}
		p.info.Duration = track.Duration
	}
	if track, ok := m.Track(false); ok {
		d, err := m.DecodeAudio(track)
		if err != nil {
			if p.video != nil {
				p.video.Close()
			}
			return nil, err
		}
		p.audio = d
		p.info.HasAudio = true
		p.info.SampleRate = track.SampleRate
		p.info.Channels = track.Channels
		if p.info.Duration == 0 {
			p.info.Duration = track.Duration
		}
		// About a quarter of a second of slack between the decoder and the
		// device, which is enough to ride out a slow frame and short enough
		// that a seek is not heard as a delay.
		p.pcm = make([]float32, max(track.SampleRate*track.Channels/4, 4096))
	}
	if p.video == nil && p.audio == nil {
		return nil, errors.New("media: the file has nothing this device can play")
	}

	if p.video != nil {
		cv, err := canvas.NewCanvas(p.info.Width, p.info.Height)
		if err != nil {
			p.Close()
			return nil, err
		}
		p.cur = cv
		if cv, err = canvas.NewCanvas(p.info.Width, p.info.Height); err == nil {
			p.next = cv
		}
	}
	return p, nil
}

// Info is what the file holds.
func (p *Player) Info() Info { return p.info }

// Play starts, or carries on after a pause.
func (p *Player) Play() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.playing {
		return
	}
	p.playing = true
	p.started = time.Now()
}

// Pause stops the clock. Nothing is thrown away.
func (p *Player) Pause() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.playing {
		return
	}
	p.base = p.position()
	p.playing = false
}

// Playing reports whether the clock is running.
func (p *Player) Playing() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.playing
}

// Position is where the file is up to.
func (p *Player) Position() time.Duration {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.position()
}

// position is Position without the lock, which the callers below already
// hold.
func (p *Player) position() time.Duration {
	if p.info.HasAudio && p.info.SampleRate > 0 && p.played > 0 {
		frames := p.played / int64(max(p.info.Channels, 1))
		return time.Duration(frames) * time.Second / time.Duration(p.info.SampleRate)
	}
	if !p.playing {
		return p.base
	}
	return p.base + time.Since(p.started)
}

// Frame is the picture that belongs on screen now, and whether there is one.
//
// The canvas belongs to the player and is handed out again and again; a
// frame that has to outlive the next call has to be copied. It reports false
// before the first frame has been decoded and after the file has ended.
func (p *Player) Frame() (*canvas.Canvas, bool, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed || p.video == nil {
		return nil, false, nil
	}
	now := p.position()

	// Decode ahead until the frame in hand is the right one for now. More
	// than a few in one go means the decoder is behind, and running the loop
	// to the end would make it later still — so it catches up a little each
	// frame instead of stalling one.
	for range 4 {
		if !p.haveNext {
			if err := p.decodeNext(); err != nil {
				return nil, false, err
			}
		}
		if !p.haveNext || p.nextPts > now {
			break
		}
		p.cur, p.next = p.next, p.cur
		p.curPts = p.nextPts
		p.haveNext = false
	}
	if p.curPts == 0 && !p.started.IsZero() && p.cur == nil {
		return nil, false, nil
	}
	return p.cur, true, nil
}

// decodeNext pulls one frame out of the decoder into the spare canvas.
func (p *Player) decodeNext() error {
	if p.next == nil {
		return nil
	}
	words := unsafe.Slice((*uint32)(unsafe.Pointer(&p.next.Pixels[0])), len(p.next.Pixels))
	w, h, pts, ok, err := p.video.Frame(words, p.next.Stride)
	if err != nil {
		return err
	}
	if p.video.Ended() {
		p.ended = true
	}
	if !ok {
		return nil
	}
	if w != p.next.Width || h != p.next.Height {
		// The decoder settled on a different size than the container
		// claimed. Make canvases that fit and lose this one frame.
		cv, err := canvas.NewCanvas(w, h)
		if err != nil {
			return err
		}
		p.next = cv
		if spare, err := canvas.NewCanvas(w, h); err == nil {
			p.cur = spare
		}
		p.info.Width, p.info.Height = w, h
		return nil
	}
	p.nextPts = pts
	p.haveNext = true
	return nil
}

// Sound is what to hand to [antui.OpenAudio]: it fills a block with the
// file's sound and keeps the clock.
//
// It is nil when the file has no sound, and a caller should check rather
// than opening a silent device.
func (p *Player) Sound() func([]float32) {
	if p.audio == nil {
		return nil
	}
	return p.fill
}

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

func (p *Player) buffered() int {
	if p.pcmTail >= p.pcmHead {
		return p.pcmTail - p.pcmHead
	}
	return len(p.pcm) - p.pcmHead + p.pcmTail
}

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

// Ended reports whether the file has played out.
func (p *Player) Ended() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.ended
}

// Seek jumps to a position. Both tracks jump to the nearest key frame at or
// before it, which for the picture may be as much as a second earlier —
// there is no other kind of seek in a compressed stream.
func (p *Player) Seek(to time.Duration) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return errors.New("media: the player is closed")
	}
	for _, d := range []*ndk.Decoder{p.video, p.audio} {
		if d == nil {
			continue
		}
		if err := d.Seek(to); err != nil {
			return err
		}
	}
	p.base = to
	p.started = time.Now()
	p.played = 0
	p.pcmHead, p.pcmTail = 0, 0
	p.haveNext = false
	p.curPts, p.nextPts = to, to
	p.ended = false
	return nil
}

// Close stops everything and gives it back.
func (p *Player) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return nil
	}
	p.closed = true
	p.playing = false
	if p.video != nil {
		p.video.Close()
	}
	if p.audio != nil {
		p.audio.Close()
	}
	p.cur, p.next = nil, nil
	return nil
}
