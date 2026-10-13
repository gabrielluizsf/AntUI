//go:build android

package antui

import "time"

func (d *androidAudio) latency() time.Duration { return d.stream.Latency() }
