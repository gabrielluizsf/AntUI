//go:build android

package audio

// Channels is how many the device is really giving.
func (r *Recorder) Channels() int { return r.stream.Channels }
