//go:build android

package audio

import "time"

// Close stops recording.
func (r *Recorder) Close() error {
	r.stopOnce.Do(func() {
		close(r.stop)
		select {
		case <-r.done:
		case <-time.After(2 * time.Second):
		}
	})
	return r.stream.Close()
}
