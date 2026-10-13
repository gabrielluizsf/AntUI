//go:build android

package ndk

import "time"

// Latency is roughly how long it is between writing a sample and hearing it:
// the buffer over the rate. It is what the device was talked into rather
// than a measurement of the whole path, which no API reports.
func (s *AudioStream) Latency() time.Duration {
	if s.Rate <= 0 {
		return 0
	}
	return time.Duration(s.BufferFrames()) * time.Second / time.Duration(s.Rate)
}
