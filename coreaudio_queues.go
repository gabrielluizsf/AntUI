//go:build darwin && cgo

package antui

import "sync"

// The queues that are open, by a number rather than by a pointer: a Go
// pointer cannot be handed to C and kept there, and this is the one thing C
// has to hold on to.
var (
	queuesMu  sync.Mutex
	queues    = map[int]*coreaudio{}
	nextQueue int
)
