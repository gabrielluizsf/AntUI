package antui

import (
	"time"
)

// A driver is one platform's way of doing it.
type driver interface {
	close() error
	latency() time.Duration
}
