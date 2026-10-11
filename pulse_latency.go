//go:build linux && !android

package antui

import (
	"time"
)

func (p *pulse) latency() time.Duration { return p.delay }
