//go:build android

package app

import (
	"unsafe"
)

type message struct {
	ev  Event
	ptr unsafe.Pointer
	// ack is closed by the app goroutine when the event has been dealt with.
	// It is nil for the events the UI thread does not wait on.
	ack chan struct{}
}
