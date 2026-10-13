//go:build android

package audio

import "time"

const recorderCloseTimeout = 2 * time.Second

// Close stops recording.
func (r *Recorder) Close() error {
	r.stopOnce.Do(func() {
		close(r.stop)
		select {
		case <-r.done:
		case <-time.After(recorderCloseTimeout):
		}
	})
	return r.stream.Close()
}
