package antui

import "time"

// Latency is roughly how long it is between a sample being written and being
// heard. It is what the device was talked into rather than a measurement.
func (a *Audio) Latency() time.Duration {
	if a == nil || a.drv == nil {
		return 0
	}
	return a.drv.latency()
}
