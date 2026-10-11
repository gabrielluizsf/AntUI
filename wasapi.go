//go:build windows

package antui

import (
	"time"
)

// A wasapi is one render stream and the thread feeding it.
type wasapi struct {
	audio *Audio

	ready chan error
	quit  chan struct{}
	done  chan struct{}

	event   uintptr
	delay   time.Duration
	stopped bool
}
