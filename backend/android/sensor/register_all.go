//go:build android

package sensor

import (
	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

// registerAll subscribes to every chosen sensor.
func registerAll(queue *ndk.SensorQueue, chosen []ndk.Sensor, rate int, latency int64) error {
	for _, sn := range chosen {
		// Batching first, and the plain path when the platform is too old
		// for it or the caller asked for none.
		if latency > 0 {
			if err := queue.Register(sn, rate, latency); err == nil {
				continue
			}
		}
		if err := queue.Enable(sn); err != nil {
			return err
		}
		if rate > 0 {
			if err := queue.SetRate(sn, rate); err != nil {
				return err
			}
		}
	}
	return nil
}
