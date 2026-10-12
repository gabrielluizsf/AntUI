//go:build android

package app

import (
	"syscall"
)

// RunOnUI puts f on the UI thread's queue and returns without waiting.
func RunOnUI(f func()) error {
	uiMu.Lock()
	if uiWrite < 0 {
		uiMu.Unlock()
		return ErrNoUIThread
	}
	uiQueue = append(uiQueue, f)
	fd := uiWrite
	uiMu.Unlock()
	// One byte is a nudge, not a message: the queue is what carries the
	// work, and a pipe that already has bytes in it wakes the looper just
	// the same.
	if _, err := syscall.Write(fd, []byte{0}); err != nil && err != syscall.EAGAIN {
		return err
	}
	return nil
}
