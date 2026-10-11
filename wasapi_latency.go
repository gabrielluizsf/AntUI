//go:build windows

package antui

import (
	"time"
)

func (w *wasapi) latency() time.Duration { return w.delay }
