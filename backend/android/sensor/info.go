//go:build android

package sensor

import (
	"time"
)

// Info describes a sensor the device has.
type Info struct {
	Kind       Kind
	Name       string
	Vendor     string
	Resolution float32
	// MinDelay is the shortest interval between readings. It is zero for a
	// sensor that reports only when something changes, and asking one of
	// those for a rate does nothing.
	MinDelay time.Duration
}
