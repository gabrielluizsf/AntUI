//go:build android

package audio

import "time"

func (r *Recorder) run(f func([]float32)) {
	defer close(r.done)

	frames := r.stream.Burst * blockBursts
	if frames <= 0 {
		frames = 1024
	}
	block := make([]float32, frames*r.stream.Channels)

	for {
		select {
		case <-r.stop:
			return
		default:
		}
		n, err := r.stream.Read(block, time.Second)
		if err != nil {
			r.mu.Lock()
			if r.err == nil {
				r.err = err
			}
			r.mu.Unlock()
			return
		}
		if n > 0 {
			f(block[:n*r.stream.Channels])
		}
	}
}
