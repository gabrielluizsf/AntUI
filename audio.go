package antui

import (
	"sync"
)

// An Audio is one stream of sound going out to the speakers.
type Audio struct {
	rate     int
	channels int
	fill     func([]float32)

	mu   sync.Mutex
	err  error
	drv  driver
	done bool
}
