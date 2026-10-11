package antui

import (
	"time"
)

// handOver passes the frame deadline to a platform that can wait for it.
func (win *Window) handOver(due time.Time) {
	if pacer, ok := win.native.(framePacer); ok {
		pacer.setFrameDeadline(due)
	}
}
