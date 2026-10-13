//go:build android

package antui

import "time"

func (d *androidAudio) run(a *Audio) {
	defer close(d.done)

	frames := d.stream.Burst * blockBursts
	if frames <= 0 {
		frames = 1024
	}
	block := make([]float32, frames*d.stream.Channels)

	for {
		select {
		case <-d.stop:
			return
		default:
		}
		a.pull(block)
		// A timeout rather than a wait for ever: a device that has been
		// taken away — a headset unplugged mid-write — otherwise leaves this
		// goroutine here for the life of the process.
		if _, err := d.stream.Write(block, time.Second); err != nil {
			if !a.stopped() {
				a.fail(err)
			}
			return
		}
	}
}
