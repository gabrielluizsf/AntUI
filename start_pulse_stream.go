//go:build linux && !android

package antui

import "time"

// startStream creates the playback stream and starts feeding it.
func (p *pulse) startStream(name string, a *Audio) error {
	// What the buffer is worth in bytes. A game wants this short — a sound
	// that arrives a tenth of a second after the thing that made it is a
	// sound that belongs to nothing — and not so short that a frame taking
	// longer than usual leaves a hole in it.
	const (
		wanted = 25 * time.Millisecond // between writing a sample and hearing it
		least  = 10 * time.Millisecond // the least the server should ask for
	)
	frame := uint32(a.channels * 4)
	target := uint32(a.rate) * frame * uint32(wanted) / uint32(time.Second)
	minimum := uint32(a.rate) * frame * uint32(least) / uint32(time.Second)

	reply, err := p.ask(cmdCreatePlaybackStream, func(t *tags) {
		writeCreateStream(t, pulsePlayback{name: name, audio: a, target: target, minimum: minimum})
	})
	if err != nil {
		return err
	}

	channel, err := reply.u32()
	if err != nil {
		return err
	}
	if _, err := reply.u32(); err != nil { // the sink input's own index
		return err
	}
	missing, err := reply.u32()
	if err != nil {
		return err
	}
	p.channel = channel
	p.delay = time.Duration(target) * time.Second / time.Duration(uint32(a.rate)*frame)

	// The block handed to the game each time. Big enough that the goroutine
	// is not woken for a handful of samples, small enough to fit inside the
	// latency asked for.
	frames := max(int(minimum/frame), 64)
	p.block = make([]float32, frames*a.channels)
	p.out = make([]byte, len(p.block)*4)

	go p.run(missing)
	return nil
}
