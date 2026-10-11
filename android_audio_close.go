//go:build android

package antui

import (
	"time"
)

func (d *androidAudio) close() error {
	d.stopOnce.Do(func() {
		close(d.stop)
		// Wait for the goroutine to stop writing before the stream goes:
		// closing one that is being written to is a use-after-free inside
		// the platform.
		select {
		case <-d.done:
		case <-time.After(2 * time.Second):
		}
	})
	return d.stream.Close()
}
